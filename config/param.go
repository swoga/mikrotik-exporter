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
		return errors.New("unknown datetime_Type")
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

type remapRe struct {
	regex       *regexp.Regexp
	replacement *string
}

func (r *remapRe) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("expected map")
	}

	if len(node.Content) != 2 {
		return errors.New("expected map with one key value pair")
	}

	var exprRaw interface{}

	if err := node.Content[0].Decode(&exprRaw); err != nil {
		return err
	}
	expr, ok := exprRaw.(string)
	if !ok {
		return errors.New("expression not string")
	}
	var err error
	r.regex, err = regexp.Compile(expr)
	if err != nil {
		return err
	}

	var replacementRaw interface{}
	if err := node.Content[1].Decode(&replacementRaw); err != nil {
		return err
	}
	if replacementRaw != nil {
		replacement, ok := replacementRaw.(string)
		if !ok {
			return errors.New("replacement not string or nil")
		}
		r.replacement = &replacement
	}

	return nil
}
