package config

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
	"github.com/swoga/go-routeros/proto"
	"github.com/swoga/mikrotik-exporter/utils"
)

type contextCommandNo struct{}

func (x *Command) Run(ctx context.Context, log zerolog.Logger, client *routeros.Client, registerer prometheus.Registerer, parentVariables map[string]string, metricCache map[string]AddMetric) error {
	commandLog := log.With().Str("command_no", ctx.Value(contextCommandNo{}).(string)).Logger()
	command := utils.Substitute(log, x.CommandBase.Command, parentVariables)
	commandLog.Debug().Str("command", command).Msg("run command")

	response, err := client.ListenArgs(strings.Split(command, "\n"))
	if err != nil {
		return err
	}

	ownCtx, cancel := context.WithTimeout(ctx, x.CommandBase.timeoutDuration)
	defer cancel()

	var i int
	for {
		select {
		case re, open := <-response.Chan():
			if !open {
				if response.Done != nil && len(response.Done.Map) > 0 {
					responseLog := commandLog.With().Str("sentence_no", "done").Logger()
					err := x.processResponse(ctx, responseLog, client, registerer, parentVariables, response.Done, metricCache)
					if err != nil {
						return err
					}
				}
				commandLog.Trace().Msg("response channel closed")
				if err := response.Err(); err != nil {
					return err
				}
			}
			if re == nil {
				if open {
					return errors.New("unexpected nil sentence while channel is open")
				}
				commandLog.Trace().Msg("all rows received")
				return nil
			}
			responseLog := commandLog.With().Int("sentence_no", i).Logger()
			i += 1
			err = x.processResponse(ctx, responseLog, client, registerer, parentVariables, re, metricCache)
			if err != nil {
				return err
			}
		case <-ownCtx.Done():
			return ownCtx.Err()
		}
	}
}

func (x *Command) processResponse(ctx context.Context, log zerolog.Logger, client *routeros.Client, registerer prometheus.Registerer, variables map[string]string, re *proto.Sentence, metricCache map[string]AddMetric) error {
	log.Trace().Interface("re", re.Map).Msg("response")

	x.addMetrics(log, registerer, variables, re, metricCache)

	childVariables := x.getChildVariables(log, re.Map, variables)
	err := x.runSubCommands(ctx, log, client, registerer, childVariables, metricCache)
	if err != nil {
		return err
	}

	return nil
}

func (x *Command) addMetrics(log zerolog.Logger, registerer prometheus.Registerer, variables map[string]string, re *proto.Sentence, metricCache map[string]AddMetric) {
	if x.Metrics == nil {
		return
	}

	commandLabelNames := x.Labels.LabelNames()
	commandLabelValues := x.Labels.LabelValues(log, re.Map, variables)

	for _, metric := range *x.Metrics {
		value, ok := metric.TryGetValue(log, re.Map, variables)
		if !ok {
			continue
		}

		metric.AddValue(log, registerer, value, commandLabelNames, commandLabelValues, re.Map, variables, metricCache)
	}
}

func (x *Command) getChildVariables(log zerolog.Logger, response map[string]string, variables map[string]string) map[string]string {
	childVariables := utils.CopyStringStringMap(variables)
	for _, variable := range x.Variables {
		childVariables[variable.LabelName] = variable.AsString(log, response, variables)
	}
	return childVariables
}

func (x *Command) runSubCommands(ctx context.Context, log zerolog.Logger, client *routeros.Client, registerer prometheus.Registerer, variables map[string]string, metricCache map[string]AddMetric) error {
	for i, subCommand := range x.SubCommands {
		commandNo := ctx.Value(contextCommandNo{}).(string)
		commandCtx := context.WithValue(ctx, contextCommandNo{}, commandNo+","+strconv.Itoa(i))

		err := subCommand.Run(commandCtx, log, client, registerer, variables, metricCache)
		if err != nil {
			return err
		}
	}

	return nil
}
