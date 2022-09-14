package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

type Credentials struct {
	Username *string `yaml:"username"`
	Password *string `yaml:"password"`
}

func (c Credentials) MarshalYAML() (interface{}, error) {
	type plain Credentials
	v := "<redacted>"
	return plain{
		Username: &v,
		Password: &v,
	}, nil
}

type Config struct {
	ConfigD ConfigD `yaml:",inline"`

	Listen      string `yaml:"listen"`
	MetricsPath string `yaml:"metrics_path"`
	ProbePath   string `yaml:"probe_path"`
	ReloadPath  string `yaml:"reload_path"`
	Namespace   string `yaml:"namespace"`

	Credentials Credentials `yaml:",inline"`

	ConfigFiles []string `yaml:"config_files"`

	ConnectionCleanupInterval         int           `yaml:"connection_cleanup_interval"`
	ConnectionCleanupIntervalDuration time.Duration `yaml:"-"`
	ConnectionUseTimeout              int           `yaml:"connection_use_timeout"`
	ConnectionUseTimeoutDuration      time.Duration `yaml:"-"`

	moduleMap map[string]*Module
	TargetMap map[string]*Target `yaml:"-"`
}

func DefaultConfig() Config {
	return Config{
		Listen:                    ":9436",
		MetricsPath:               "/metrics",
		ProbePath:                 "/probe",
		ReloadPath:                "/-/reload",
		Namespace:                 "mikrotik",
		ConfigFiles:               []string{"./conf.d/*"},
		ConnectionCleanupInterval: 60,
		ConnectionUseTimeout:      300,

		moduleMap: make(map[string]*Module),
		TargetMap: make(map[string]*Target),
	}
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig()

	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	c.ConnectionCleanupIntervalDuration = time.Duration(c.ConnectionCleanupInterval) * time.Second
	c.ConnectionUseTimeoutDuration = time.Duration(c.ConnectionUseTimeout) * time.Second

	return nil
}

func (c *Config) loadContents(basePath string) error {
	if err := c.loadConfigFiles(basePath); err != nil {
		return err
	}

	if err := c.populateModuleMap(); err != nil {
		return err
	}

	if err := c.populateTargetMap(); err != nil {
		return err
	}

	if err := c.populateTargetCredentials(); err != nil {
		return err
	}

	c.applyExtensions()

	return nil
}

func (c *Config) loadConfigFiles(basePath string) error {
	err := os.Chdir(basePath)
	if err != nil {
		return err
	}

	for _, path := range c.ConfigFiles {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

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
	decoder := yaml.NewDecoder(yamlReader, yaml.Strict())

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
		log.Logger.Trace().Str("module", module.Name).Msg("add module")
		c.moduleMap[module.Name] = module
	}
	return nil
}

func (c *Config) populateTargetMap() error {
	log.Logger.Trace().Msg("populate target map")
	for _, target := range c.ConfigD.Targets {
		_, exists := c.TargetMap[target.Name]
		if exists {
			return fmt.Errorf("non-unique target name: %s", target.Name)
		}
		log.Logger.Trace().Str("target", target.Name).Msg("add target")
		c.TargetMap[target.Name] = target
	}
	return nil
}

func (c *Config) populateTargetCredentials() error {
	log.Logger.Trace().Msg("populate target credentials")

	for _, target := range c.TargetMap {
		if target.Credentials.Username == nil {
			target.Credentials.Username = c.Credentials.Username
			log.Logger.Trace().Str("target", target.Name).Msg("use global username")
		}
		if target.Credentials.Password == nil {
			target.Credentials.Password = c.Credentials.Password
			log.Logger.Trace().Str("target", target.Name).Msg("use global password")
		}
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
	target, ok := c.TargetMap[name]
	if !ok {
		return nil
	}
	return target
}
