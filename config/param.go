package config

import (
	"errors"
	"regexp"

	"github.com/swoga/mikrotik-exporter/utils"
	"gopkg.in/yaml.v3"
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

	RemapValues        map[string]*string `yaml:"remap_values"`
	RemapValuesRe      map[string]*string `yaml:"remap_values_re"`
	remapValueReParsed map[*regexp.Regexp]*string

	Negate       bool   `yaml:"negate"`
	DateTimeType string `yaml:"datetime_type"`
}

func DefaultParam() Param {
	return Param{
		remapValueReParsed: make(map[*regexp.Regexp]*string),
	}
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
		return errors.New("unknown datetime_Type")
	}
	for regex, value := range param.RemapValuesRe {
		expr, err := regexp.Compile(regex)
		if err != nil {
			return err
		}
		param.remapValueReParsed[expr] = value
	}

	return nil
}

func (param *Param) UnmarshalYAML(node *yaml.Node) error {
	*param = DefaultParam()

	type plain Param
	if err := node.Decode((*plain)(param)); err != nil {
		return err
	}

	err := param.Validate()
	if err != nil {
		return err
	}

	return nil
}
