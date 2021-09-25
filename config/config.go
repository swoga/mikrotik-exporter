package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	ConfigD ConfigD `yaml:",inline"`

	Listen      string `yaml:"listen"`
	MetricsPath string `yaml:"metrics_path"`
	ProbePath   string `yaml:"probe_path"`
	ReloadPath  string `yaml:"reload_path"`

	Credentials Credentials `yaml:",inline"`

	ConfigFiles []string `yaml:"config_files"`

	moduleMap map[string]*Module
	targetMap map[string]*Target
}

func DefaultConfig() Config {
	return Config{
		Listen:      ":9436",
		MetricsPath: "/metrics",
		ProbePath:   "/probe",
		ReloadPath:  "/-/reload",
		ConfigFiles: []string{"./conf.d/*"},
		moduleMap:   make(map[string]*Module),
		targetMap:   make(map[string]*Target),
	}
}

func (c *Config) UnmarshalYAML(node *yaml.Node) error {
	*c = DefaultConfig()

	type plain Config
	if err := node.Decode((*plain)(c)); err != nil {
		return err
	}

	if err := c.loadConfigFiles(); err != nil {
		return err
	}

	if err := c.populateModuleMap(); err != nil {
		return err
	}

	if err := c.populateTargetMap(); err != nil {
		return err
	}

	c.applyExtensions()

	return nil
}

func (c *Config) loadConfigFiles() error {
	for _, path := range c.ConfigFiles {
		log.Logger.Debug().Str("path", path).Msg("get files for glob")
		files, err := filepath.Glob(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = c.loadConfigFile(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) loadConfigFile(configFile string) error {
	log.Logger.Info().Str("file", configFile).Msg("load sub-config")

	yamlReader, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("error reading sub-config file: %s", err)
	}
	defer yamlReader.Close()
	decoder := yaml.NewDecoder(yamlReader)
	decoder.KnownFields(true)

	data := &ConfigD{}
	err = decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error parsing sub-config file: %s", err)
	}

	c.mergeConfig(*data)

	return nil
}

func (c *Config) mergeConfig(data ConfigD) {
	log.Logger.Trace().Msg("merge sub-config data")
	c.ConfigD.Modules = append(c.ConfigD.Modules, data.Modules...)
	c.ConfigD.ModuleExtensions = append(c.ConfigD.ModuleExtensions, data.ModuleExtensions...)
}

func (c *Config) populateModuleMap() error {
	log.Logger.Trace().Msg("populate module map")
	for _, module := range c.ConfigD.Modules {
		_, exists := c.moduleMap[module.Name]
		if exists {
			return fmt.Errorf("non-unique module name: %s", module.Name)
		}
		c.moduleMap[module.Name] = module
	}
	return nil
}

func (c *Config) populateTargetMap() error {
	log.Logger.Trace().Msg("populate target map")
	for _, target := range c.ConfigD.Targets {
		_, exists := c.targetMap[target.Name]
		if exists {
			return fmt.Errorf("non-unique target name: %s", target.Name)
		}
		c.targetMap[target.Name] = target
	}
	return nil
}

func (c *Config) GetModule(name string) *Module {
	module, ok := c.moduleMap[name]
	if !ok {
		return nil
	}
	return module
}

func (c *Config) applyExtensions() {
	log.Logger.Trace().Msg("apply extensions")
	for _, ext := range c.ConfigD.ModuleExtensions {
		module := c.GetModule(ext.Name)
		if module == nil {
			log.Logger.Warn().Str("module", ext.Name).Msg("try to extend non-existing module")
			continue
		}

		ext.ExtendModule(log.Logger, *module)
	}
	c.ConfigD.ModuleExtensions = nil
}

func (c *Config) GetTarget(name string) *Target {
	target, ok := c.targetMap[name]
	if !ok {
		return nil
	}
	return target
}
