package config

import (
	"errors"
	"strings"
	"time"
)

type Target struct {
	Name            string `yaml:"name"`
	Address         string `yaml:"address"`
	Timeout         int    `yaml:"timeout"`
	TimeoutDuration time.Duration
	Queue           int               `yaml:"queue"`
	Credentials     Credentials       `yaml:",inline"`
	Variables       map[string]string `yaml:"variables"`
	Modules         []string          `yaml:"modules"`
}

func DefaultTarget() Target {
	return Target{
		Timeout: 10,
		Queue:   1000,
	}
}

func (target *Target) Validate() error {
	if target.Name == "" {
		return errors.New("require name")
	}
	if target.Address == "" {
		return errors.New("require address")
	}
	if !strings.Contains(target.Address, ":") {
		target.Address = target.Address + ":8728"
	}

	return nil
}

func (target *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*target = DefaultTarget()

	type plain Target
	if err := unmarshal((*plain)(target)); err != nil {
		return err
	}

	target.TimeoutDuration = time.Duration(target.Timeout) * time.Second

	err := target.Validate()
	if err != nil {
		return err
	}

	return nil
}
