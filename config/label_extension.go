package config

type LabelExtension struct {
	Extension Extension `yaml:",inline"`
	Label     Label     `yaml:",inline"`
}
