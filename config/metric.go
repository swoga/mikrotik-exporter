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
	Help       string `yaml:"help"`
	Labels     Labels `yaml:"labels"`
}

func DefaultMetric() Metric {
	return Metric{}
}

func (x *Metric) GetName() string {
	if x.MetricName != "" {
		return x.MetricName
	}
	return strings.ReplaceAll(x.Param.ParamName, "-", "_")
}

func (metric *Metric) Validate() error {
	if metric.MetricType == "" {
		return errors.New("require metric_type")
	}
	if !utils.ArrayContainsString(metricTypes, metric.MetricType) {
		return errors.New("unknown metric_type")
	}
	if metric.GetName() == "" {
		return errors.New("require metric_name or param_name")
	}
	if !regexValidMetricAndLabelName.MatchString(metric.GetName()) {
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
