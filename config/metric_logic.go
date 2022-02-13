package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func (metric *Metric) TryGetValue(log zerolog.Logger, response map[string]string, variables map[string]string) (float64, bool) {
	paramLog := log.With().Str("metric_name", metric.GetName()).Logger()
	paramLog.Trace().Msg("get metric value")
	value, ok := metric.Param.tryGetValue(paramLog, response, variables)
	if ok {
		paramLog.Trace().Float64("value", value).Msg("got metric value")
	}
	return value, ok
}

type AddMetric func(labelValues []string, value float64)

func (metric *Metric) createPrometheusGauge(registerer prometheus.Registerer, labelNames []string) (AddMetric, error) {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metric.GetName(),
		Help: metric.Help,
	}, labelNames)
	err := registerer.Register(vec)

	if err != nil {
		return nil, err
	}

	fn := func(labelValues []string, value float64) {
		vec.WithLabelValues(labelValues...).Set(value)
	}
	return fn, nil
}

func (metric *Metric) createPrometheusCounter(registerer prometheus.Registerer, labelNames []string) (AddMetric, error) {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metric.GetName(),
		Help: metric.Help,
	}, labelNames)
	err := registerer.Register(vec)

	if err != nil {
		return nil, err
	}

	fn := func(labelValues []string, value float64) {
		vec.WithLabelValues(labelValues...).Add(value)
	}
	return fn, nil
}

func (metric *Metric) createPrometheusMetric(log zerolog.Logger, registerer prometheus.Registerer, commandLabelNames []string) (AddMetric, error) {
	metricLabelNames := metric.HasLabels.LabelNames()
	labelNames := append(commandLabelNames, metricLabelNames...)

	switch metric.MetricType {
	case METRIC_TYPE_GAUGE:
		return metric.createPrometheusGauge(registerer, labelNames)
	case METRIC_TYPE_COUNTER:
		return metric.createPrometheusCounter(registerer, labelNames)
	default:
		log.Panic().Str("metric_type", metric.MetricType).Msg("invalid metric_type")
		return nil, nil
	}
}

func (metric *Metric) AddValue(log zerolog.Logger, registerer prometheus.Registerer, value float64, commandLabelNames, commandLabelValues []string, response map[string]string, variables map[string]string, metricCache map[string]AddMetric) {
	metricLabelValues := metric.HasLabels.LabelValues(log, response, variables)
	labelValues := append(commandLabelValues, metricLabelValues...)

	fn, found := metricCache[metric.GetName()]
	if !found {
		fnNew, err := metric.createPrometheusMetric(log, registerer, commandLabelNames)
		if err != nil {
			panic(err)
		}
		metricCache[metric.GetName()] = fnNew
		fn = fnNew
	}

	fn(labelValues, value)
}
