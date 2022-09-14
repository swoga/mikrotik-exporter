package config

type ConfigD struct {
	Targets          []*Target          `yaml:"targets"`
	Modules          []*Module          `yaml:"modules"`
	ModuleExtensions []*ModuleExtension `yaml:"module_extensions"`
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
