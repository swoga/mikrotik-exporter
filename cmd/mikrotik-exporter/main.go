package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swoga/mikrotik-exporter/config"
	"github.com/swoga/mikrotik-exporter/connection"
	"go.uber.org/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version string = "dev"

	requests          atomic.Int64
	sc                *config.SafeConfig
	connectionManager *connection.ConnectionManager
	consoleWriter     = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
)

func main() {
	log.Logger = log.Output(consoleWriter)
	log.Info().Str("version", version).Msg("starting mikrotik-exporter")

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
	c := sc.Get()
	connectionManager = connection.CreateConnectionManager(c.ConnectionCleanupIntervalDuration, c.ConnectionUseTimeoutDuration)

	log.Info().Str("path", c.MetricsPath).Msg("serve internal metrics at")
	http.Handle(c.MetricsPath, promhttp.Handler())
	log.Info().Str("path", c.ProbePath).Msg("listen for probe requests at")
	http.HandleFunc(c.ProbePath, handleProbeRequest)
	log.Info().Str("listen", c.Listen).Msg("starting http server")
	server := &http.Server{Addr: c.Listen}
	go startServer(server)

	waitForSigterm()

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(shutdownContext)
}

func loadConfig(configFile string) *config.SafeConfig {
	loader, err := config.New(configFile)
	if err == nil {
		err = loader.LoadConfig()
	}
	if err != nil {
		log.Panic().Err(err).Msg("error loading config file")
	}

	loader.EnableReloadByHTTP()
	loader.EnableReloadBySIGHUP()

	return loader
}

func startServer(server *http.Server) {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("error starting http server")
	}
}

func waitForSigterm() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	<-term
	log.Info().Msg("received SIGTERM, shutting down")
}
