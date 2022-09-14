package config

import (
	"errors"
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

func (module *Module) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*module = DefaultModule()

	type plain Module
	if err := unmarshal((*plain)(module)); err != nil {
		return err
	}

	err := module.Validate()
	if err != nil {
		return err
	}

	return nil
}
