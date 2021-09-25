package config

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Target struct {
	Name            string `yaml:"name"`
	Address         string `yaml:"address"`
	Timeout         int    `yaml:"timeout"`
	timeoutDuration time.Duration
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

func (target *Target) UnmarshalYAML(node *yaml.Node) error {
	*target = DefaultTarget()

	type plain Target
	if err := node.Decode((*plain)(target)); err != nil {
		return err
	}

	target.timeoutDuration = time.Duration(target.Timeout) * time.Second

	err := target.Validate()
	if err != nil {
		return err
	}

	return nil
}
