package config

import "github.com/rs/zerolog"

type ModuleExtension struct {
	Name     string            `yaml:"name"`
	Commands CommandExtensions `yaml:"commands"`
}

func (x *ModuleExtension) ExtendModule(log zerolog.Logger, module Module) {
	moduleExtLog := log.With().Str("module", module.Name).Logger()
	moduleExtLog.Trace().Msg("extend module")
	x.Commands.Extend(moduleExtLog, module.Commands)
}
