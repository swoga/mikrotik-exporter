package config

import (
	"errors"
	"strings"

	"github.com/swoga/mikrotik-exporter/utils"
)

const (
	METRIC_TYPE_COUNTER = "counter"
	METRIC_TYPE_GAUGE   = "gauge"
)

var (
	metricTypes = []string{METRIC_TYPE_COUNTER, METRIC_TYPE_GAUGE}
)

type Metric struct {
	Param Param `yaml:",inline"`

	MetricName string `yaml:"metric_name"`
	MetricType string `yaml:"metric_type"`
	Help       string `yaml:"help,omitempty"`
	Labels     Labels `yaml:"labels,omitempty"`
}

func DefaultMetric() Metric {
	return Metric{}
}

func (metric *Metric) Validate() error {
	if metric.MetricType == "" {
		return errors.New("require metric_type")
	}
	if !utils.ArrayContainsString(metricTypes, metric.MetricType) {
		return errors.New("unknown metric_type")
	}
	if metric.MetricName == "" {
		metric.MetricName = strings.ReplaceAll(metric.Param.ParamName, "-", "_")
	}
	if metric.MetricName == "" {
		return errors.New("require metric_name or param_name")
	}
	if !regexValidMetricAndLabelName.MatchString(metric.MetricName) {
		return errors.New("invalid metric_name")
	}

	return nil
}

func (x *Metric) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*x = DefaultMetric()

	type plain Metric
	if err := unmarshal((*plain)(x)); err != nil {
		return err
	}

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}
