package config

type ConfigD struct {
	Targets          []*Target          `yaml:"targets,omitempty"`
	Modules          []*Module          `yaml:"modules,omitempty"`
	ModuleExtensions []*ModuleExtension `yaml:"module_extensions,omitempty"`
}

func DefaultConfigD() ConfigD {
	return ConfigD{}
}

func (x *ConfigD) Validate() error {

	return nil
}

func (x *ConfigD) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*x = DefaultConfigD()

	type plain ConfigD
	if err := unmarshal((*plain)(x)); err != nil {
		return err
	}

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}
