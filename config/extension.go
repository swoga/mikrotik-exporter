package config

import (
	"errors"
)

const (
	EXTENSION_ACTION_ADD       = "add"
	EXTENSION_ACTION_OVERWRITE = "overwrite"
	EXTENSION_ACTION_REMOVE    = "remove"
)

type Extension struct {
	ExtensionAction string `yaml:"extension_action"`
}

func DefaultExtension() Extension {
	return Extension{}
}

func (x *Extension) Validate() error {
	if x.ExtensionAction == "" {
		return errors.New("require ExtensionAction")
	}

	return nil
}

func (x *Extension) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*x = DefaultExtension()

	type plain Extension
	if err := unmarshal((*plain)(x)); err != nil {
		return err
	}

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}
