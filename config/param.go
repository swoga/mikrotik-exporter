package config

import (
	"errors"
	"regexp"

	"github.com/swoga/mikrotik-exporter/utils"
)

const (
	PARAM_TYPE_STRING   = "string"
	PARAM_TYPE_INT      = "int"
	PARAM_TYPE_BOOL     = "bool"
	PARAM_TYPE_TIMESPAN = "timespan"
	PARAM_TYPE_DATETIME = "datetime"

	PARAM_DATETYPE_TO_NOW    = "tonow"
	PARAM_DATETYPE_FROM_NOW  = "fromnow"
	PARAM_DATETYPE_TIMESTAMP = "timestamp"
)

var (
	paramTypes         = []string{PARAM_TYPE_STRING, PARAM_TYPE_INT, PARAM_TYPE_BOOL, PARAM_TYPE_TIMESPAN, PARAM_TYPE_DATETIME}
	paramDateTimeTypes = []string{PARAM_DATETYPE_TO_NOW, PARAM_DATETYPE_FROM_NOW, PARAM_DATETYPE_TIMESTAMP}
)

type Param struct {
	ParamName string `yaml:"param_name"`
	ParamType string `yaml:"param_type"`
	Default   string `yaml:"default"`
	Value     string `yaml:"value"`

	RemapValues   map[string]*string `yaml:"remap_values"`
	RemapValuesRe []remapRe          `yaml:"remap_values_re"`

	Negate       bool   `yaml:"negate"`
	DateTimeType string `yaml:"datetime_type"`
}

func DefaultParam() Param {
	return Param{}
}

func (param *Param) Validate() error {
	if param.ParamName == "" && param.Value == "" {
		return errors.New("either param_name or value must be set")
	}

	utils.SetDefaultString(&param.ParamType, PARAM_TYPE_INT)
	if !utils.ArrayContainsString(paramTypes, param.ParamType) {
		return errors.New("unknown param_type")
	}
	utils.SetDefaultString(&param.DateTimeType, PARAM_DATETYPE_FROM_NOW)
	if !utils.ArrayContainsString(paramDateTimeTypes, param.DateTimeType) {
		return errors.New("unknown datetime_type")
	}

	return nil
}

func (param *Param) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*param = DefaultParam()

	type plain Param
	if err := unmarshal((*plain)(param)); err != nil {
		return err
	}

	err := param.Validate()
	if err != nil {
		return err
	}

	return nil
}

type remapRe struct {
	regex       *regexp.Regexp
	replacement *string
}

func (r *remapRe) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := map[string]*string{}

	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	if len(raw) != 1 {
		return errors.New("expected map with one key value pair")
	}

	for expr, replacement := range raw {
		r.regex, err = regexp.Compile(expr)
		if err != nil {
			return err
		}
		r.replacement = replacement
	}

	return nil
}
