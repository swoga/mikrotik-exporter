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

func (metric *Metric) AddValue(log zerolog.Logger, registerer prometheus.Registerer, value float64, commandLabelNames, commandLabelValues []string, response map[string]string, variables map[string]string) {
	metricLabelNames := metric.HasLabels.LabelNames()
	labelNames := append(commandLabelNames, metricLabelNames...)
	metricLabelValues := metric.HasLabels.LabelValues(log, response, variables)

	labelValues := append(commandLabelValues, metricLabelValues...)

	switch metric.MetricType {
	case METRIC_TYPE_GAUGE:
		vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metric.GetName(),
			Help: metric.Help,
		}, labelNames)
		err := registerer.Register(vec)
		if err != nil {
			switch err.(type) {
			case prometheus.AlreadyRegisteredError:
				are := err.(prometheus.AlreadyRegisteredError)
				vec = are.ExistingCollector.(*prometheus.GaugeVec)
			default:
				panic(err)
			}
		}

		vec.WithLabelValues(labelValues...).Set(value)
	case METRIC_TYPE_COUNTER:
		vec := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: metric.GetName(),
			Help: metric.Help,
		}, labelNames)
		err := registerer.Register(vec)
		if err != nil {
			switch err.(type) {
			case prometheus.AlreadyRegisteredError:
				are := err.(prometheus.AlreadyRegisteredError)
				vec = are.ExistingCollector.(*prometheus.CounterVec)
			default:
				panic(err)
			}
		}

		vec.WithLabelValues(labelValues...).Add(value)
	default:
		log.Panic().Str("metric_type", metric.MetricType).Msg("invalid metric_type")
	}
}
