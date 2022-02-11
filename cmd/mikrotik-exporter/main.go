package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swoga/mikrotik-exporter/config"
	"github.com/swoga/mikrotik-exporter/connection"
	"go.uber.org/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests          atomic.Int64
	sc                *config.SafeConfig
	connectionManager *connection.ConnectionManager
	consoleWriter     = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
)

func main() {
	log.Logger = log.Output(consoleWriter)
	configFile := flag.String("config.file", "config.yml", "")
	debug := flag.Bool("debug", false, "")
	trace := flag.Bool("trace", false, "")
	flag.Parse()

	if *trace {
		log.Logger = log.Logger.Level(zerolog.TraceLevel)
	} else if *debug {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}

	sc = loadConfig(*configFile)
	connectionManager = connection.CreateConnectionManager()

	c := sc.Get()

	log.Info().Str("path", c.MetricsPath).Msg("serve internal metrics at")
	http.Handle(c.MetricsPath, promhttp.Handler())
	log.Info().Str("path", c.ProbePath).Msg("listen for probe requests at")
	http.HandleFunc(c.ProbePath, handleProbeRequest)
	log.Info().Str("listen", c.Listen).Msg("starting http server")
	err := http.ListenAndServe(c.Listen, nil)
	if err != nil {
		log.Panic().Err(err).Msg("error starting http server")
	}
}

func loadConfig(configFile string) *config.SafeConfig {
	loader := config.New(configFile)
	err := loader.LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("error loading config file")
	}

	loader.EnableReloadByHTTP()
	loader.EnableReloadBySIGHUP()

	return &loader
}
