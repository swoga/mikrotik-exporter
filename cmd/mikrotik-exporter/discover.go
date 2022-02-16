package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Discover struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels"`
}

func handleDiscoverRequest(w http.ResponseWriter, r *http.Request) {
	requestLog := log.With().Str("request", "discover").Logger()

	c := sc.Get()
	var response []Discover

	requestLog.Trace().Msg("iterate targets")

	for name, target := range c.TargetMap {
		discover := Discover{
			Targets: []string{name},
			Labels:  target.DiscoverLabels,
		}
		response = append(response, discover)
		requestLog.Trace().Str("target", name).Interface("discover", discover).Msg("add discover for target")
	}

	b, err := yaml.Marshal(response)
	if err != nil {
		requestLog.Err(err).Msg("error during discover request")
		http.Error(w, "error during discover request", http.StatusInternalServerError)
	}

	w.Write(b)
}
