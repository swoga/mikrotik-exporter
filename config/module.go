package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Module struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

func DefaultModule() Module {
	return Module{}
}

func (module *Module) Validate() error {
	if module.Name == "" {
		return errors.New("require name")
	}

	return nil
}

func (module *Module) UnmarshalYAML(node *yaml.Node) error {
	*module = DefaultModule()

	type plain Module
	if err := node.Decode((*plain)(module)); err != nil {
		return err
	}

	err := module.Validate()
	if err != nil {
		return err
	}

	return nil
}
