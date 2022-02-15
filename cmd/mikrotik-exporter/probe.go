package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swoga/go-routeros"
	"github.com/swoga/mikrotik-exporter/config"
)

func handleProbeRequest(w http.ResponseWriter, r *http.Request) {
	debugStr := r.URL.Query().Get("debug")
	debug := false
	if debugStr != "" {
		debug = true
	}
	traceStr := r.URL.Query().Get("trace")
	trace := false
	if traceStr != "" {
		trace = true
	}

	request := requests.Inc()
	requestLog := log.With().Int64("request_no", request).Logger()

	if debug || trace {
		debugWriter := zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339, NoColor: true}
		multi := zerolog.MultiLevelWriter(consoleWriter, debugWriter)
		requestLog = requestLog.Output(multi)
		w.Header().Set("Content-Type", "text/plain")
	}

	if trace {
		requestLog = requestLog.Level(zerolog.TraceLevel)
	} else if debug {
		requestLog = requestLog.Level(zerolog.DebugLevel)
	}

	targetName := r.URL.Query().Get("target")
	if targetName == "" {
		log.Error().Msg("request with missing target")
		http.Error(w, "?target= missing", http.StatusBadRequest)
		return
	}

	requestLog = requestLog.With().Str("target", targetName).Logger()
	requestLog.Trace().Msg("received request")

	c := sc.Get()
	target := c.GetTarget(targetName)
	if target == nil {
		requestLog.Error().Msg("invalid target")
		http.Error(w, "invalid target", http.StatusNotFound)
		return
	}
	requestLog.Trace().Msg("found target")

	moduleNames := target.Modules
	queryModules := r.URL.Query().Get("modules")
	if queryModules != "" {
		requestLog.Debug().Msg("overwrite modules by query string")
		moduleNames = strings.Split(queryModules, ",")
	}

	registry := prometheus.NewRegistry()
	probeTarget(r.Context(), requestLog, c, registry, target, moduleNames)

	if debug || trace {
		mfs, _ := registry.Gather()
		for _, mf := range mfs {
			expfmt.MetricFamilyToText(w, mf)
		}
		return
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func probeTarget(ctx context.Context, log zerolog.Logger, c *config.Config, registry *prometheus.Registry, target *config.Target, moduleNames []string) {
	registerer := prometheus.WrapRegistererWithPrefix("mikrotik_exporter_", registry)

	probeSuccess := prometheus.NewGauge(prometheus.GaugeOpts{Name: "probe_success"})
	registry.Register(probeSuccess)

	start := time.Now()

	conn, err := connectionManager.Get(log, target)
	if err != nil {
		log.Err(err).Msg("error connecting to target")
		probeSuccess.Set(0)
	} else {
		err = runModules(ctx, log, c, registerer, conn.Client, moduleNames, target.Variables)
		if err != nil {
			log.Err(err).Msg("error during probe")
			probeSuccess.Set(0)
		} else {
			probeSuccess.Set(1)
		}
	}
	defer conn.Free(log, target.TimeoutDuration)

	duration := time.Since(start)
	probeDuration := prometheus.NewGauge(prometheus.GaugeOpts{Name: "probe_duration_seconds"})
	registry.Register(probeDuration)
	probeDuration.Set(duration.Seconds())
}

func runModules(ctx context.Context, log zerolog.Logger, c *config.Config, registerer prometheus.Registerer, client *routeros.Client, moduleNames []string, variables map[string]string) error {
	metricCache := make(map[string]config.AddMetric)

	for _, moduleName := range moduleNames {
		moduleLog := log.With().Str("module", moduleName).Logger()
		module := c.GetModule(moduleName)
		if module == nil {
			moduleLog.Warn().Msg("skip invalid module")
			continue
		}
		err := module.Run(ctx, moduleLog, client, registerer, variables, metricCache)
		if err != nil {
			return err
		}
	}

	return nil
}
