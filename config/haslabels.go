package config

import "github.com/rs/zerolog"

type HasLabels struct {
	Labels Labels `yaml:"labels"`
}

func (x *HasLabels) LabelNames() []string {
	names := []string{}
	if x.Labels == nil {
		return names
	}
	for _, label := range *x.Labels {
		names = append(names, label.GetName())
	}
	return names
}

func (x *HasLabels) LabelValues(log zerolog.Logger, response map[string]string, variables map[string]string) []string {
	values := []string{}
	if x.Labels == nil {
		return values
	}
	for _, label := range *x.Labels {
		labelValue := label.AsString(log, response, variables)
		values = append(values, labelValue)
	}
	return values
}
