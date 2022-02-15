package config

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
)

func (module *Module) Run(ctx context.Context, log zerolog.Logger, client *routeros.Client, registerer prometheus.Registerer, variables map[string]string, metricCache map[string]AddMetric) error {
	log.Trace().Msg("running module")
	moduleRegisterer := prometheus.WrapRegistererWithPrefix(module.Name+"_", registerer)

	for i, command := range module.Commands {
		commandCtx := context.WithValue(ctx, contextCommandNo{}, strconv.Itoa(i))
		err := command.Run(commandCtx, log, client, moduleRegisterer, variables, metricCache)
		if err != nil {
			return err
		}
	}

	return nil
}
