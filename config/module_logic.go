package config

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
)

func (module *Module) Run(ctx context.Context, log zerolog.Logger, client *routeros.Client, registerer prometheus.Registerer, variables map[string]string) error {
	log.Trace().Msg("running module")
	moduleRegisterer := prometheus.WrapRegistererWithPrefix(module.Namespace+"_", registerer)
	metricCache := make(map[string]AddMetric)

	for i, command := range module.Commands {
		commandRegisterer := moduleRegisterer
		if command.CommandBase.Prefix != "" {
			commandRegisterer = prometheus.WrapRegistererWithPrefix(command.CommandBase.Prefix+"_", moduleRegisterer)
		}
		commandCtx := context.WithValue(ctx, contextCommandNo{}, strconv.Itoa(i))
		err := command.Run(commandCtx, log, client, commandRegisterer, variables, metricCache)
		if err != nil {
			return err
		}
	}

	return nil
}
