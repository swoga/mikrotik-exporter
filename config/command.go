package config

import (
	"errors"
	"time"
)

type CommandBase struct {
	Command         string `yaml:"command"`
	Timeout         int    `yaml:"timeout"`
	timeoutDuration time.Duration
	Prefix          string `yaml:"prefix"`
}

func DefaultCommandBase() CommandBase {
	return CommandBase{
		Timeout: 10,
	}
}

func (x *CommandBase) Validate() error {
	if x.Command == "" {
		return errors.New("require command")
	}

	return nil
}

func (x *CommandBase) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*x = DefaultCommandBase()

	type plain CommandBase
	if err := unmarshal((*plain)(x)); err != nil {
		return err
	}

	x.timeoutDuration = time.Duration(x.Timeout) * time.Second

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}

type Command struct {
	CommandBase CommandBase `yaml:",inline"`

	Metrics     Metrics   `yaml:"metrics"`
	Labels      Labels    `yaml:"labels"`
	Variables   Labels    `yaml:"variables"`
	SubCommands []Command `yaml:"sub_commands"`
}

func DefaultCommand() Command {
	return Command{}
}

func (x *Command) Validate() error {

	return nil
}

func (x *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*x = DefaultCommand()

	type plain Command
	if err := unmarshal((*plain)(x)); err != nil {
		return err
	}

	err := x.Validate()
	if err != nil {
		return err
	}

	return nil
}
