package config

import "gopkg.in/yaml.v3"

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

func (x *ConfigD) UnmarshalYAML(node *yaml.Node) error {
	*x = DefaultConfigD()

	type plain ConfigD
	if err := node.Decode((*plain)(x)); err != nil {
		return err
	}

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}
