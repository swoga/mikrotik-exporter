package config

type MetricExtension struct {
	Extension `yaml:",inline"`
	Metric    `yaml:",inline"`
}
