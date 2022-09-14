package config

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	configReloadSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "mikrotik_exporter",
		Name:      "config_last_reload_successful",
	})

	configReloadSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "mikrotik_exporter",
		Name:      "config_last_reload_success_timestamp_seconds",
	})
)

func init() {
	prometheus.MustRegister(configReloadSuccess)
	prometheus.MustRegister(configReloadSeconds)
}

type SafeConfig struct {
	sync.RWMutex
	configFile string
	c          *Config
}

func (sc *SafeConfig) Get() *Config {
	sc.Lock()
	defer sc.Unlock()
	return sc.c
}

func New(configFile string) (*SafeConfig, error) {
	configFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}

	return &SafeConfig{
		c:          &Config{},
		configFile: configFileAbs,
	}, nil
}

func (sc *SafeConfig) LoadConfig() (err error) {
	c := &Config{}
	defer func() {
		if err != nil {
			configReloadSuccess.Set(0)
		} else {
			configReloadSuccess.Set(1)
			configReloadSeconds.SetToCurrentTime()
		}
	}()

	log.Logger.Info().Str("file", sc.configFile).Msg("load config")

	yamlReader, err := os.Open(sc.configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}
	defer yamlReader.Close()
	decoder := yaml.NewDecoder(yamlReader)
	decoder.KnownFields(true)

	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("error parsing config file: %s", err)
	}
	basePath := filepath.Dir(sc.configFile)
	err = c.loadContents(basePath)
	if err != nil {
		return fmt.Errorf("error loading config file: %s", err)
	}

	sc.Lock()
	sc.c = c
	defer sc.Unlock()

	return nil
}

func (sc *SafeConfig) EnableReloadByHTTP() {
	reloadRequest := make(chan chan error)
	go func() {
		for {
			reloadResult := <-reloadRequest
			log.Debug().Msg("config reload triggerd by API")
			err := sc.LoadConfig()
			reloadResult <- err

			if err != nil {
				log.Logger.Err(err).Msg("error reloading config")
			} else {
				log.Logger.Info().Msg("config reloaded")
			}
		}
	}()

	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		reloadResult := make(chan error)
		reloadRequest <- reloadResult
		err := <-reloadResult
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to reload config: %s", err), http.StatusInternalServerError)
		}
	})
}

func (sc *SafeConfig) EnableReloadBySIGHUP() {
	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	go func() {
		for {
			<-hup
			log.Debug().Msg("config reload triggerd by SIGHUP")
			err := sc.LoadConfig()

			if err != nil {
				log.Logger.Err(err).Msg("error reloading config")
			} else {
				log.Logger.Info().Msg("config reloaded")
			}
		}
	}()
}
