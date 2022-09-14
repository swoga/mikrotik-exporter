package config

import "github.com/rs/zerolog"

type ModuleExtension struct {
	Name     string             `yaml:"name"`
	Commands []CommandExtension `yaml:"commands"`
}

func (x *ModuleExtension) GetByCommand(commandStr string) []CommandExtension {
	commands := []CommandExtension{}
	for _, command := range x.Commands {
		if command.Command == commandStr {
			commands = append(commands, command)
		}
	}
	return commands
}

func (x *ModuleExtension) ExtendModule(log zerolog.Logger, module Module) {
	moduleExtLog := log.With().Str("module", module.Name).Logger()
	moduleExtLog.Trace().Msg("extend module")
	for _, command := range module.Commands {
		extensions := x.GetByCommand(command.CommandBase.Command)
		moduleExtLog.Trace().Str("command", command.CommandBase.Command).Int("n", len(extensions)).Msg("got extensions for command")
		for _, extension := range extensions {
			extension.ExtendCommand(moduleExtLog, command)
		}
	}
}
