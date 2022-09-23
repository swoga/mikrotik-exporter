package config

import (
	"errors"
	"regexp"
	"strings"
)

type Label struct {
	Param Param `yaml:",inline"`

	LabelName string `yaml:"label_name"`
}

func DefaultLabel() Label {
	return Label{}
}

var regexValidMetricAndLabelName = regexp.MustCompile("^[a-zA-Z_:][a-zA-Z0-9_:]*$")

func (label *Label) Validate() error {
	if label.LabelName == "" {
		label.LabelName = strings.ReplaceAll(label.Param.ParamName, "-", "_")
	}
	if label.LabelName == "" {
		return errors.New("require param_name or label_name")
	}
	if !regexValidMetricAndLabelName.MatchString(label.LabelName) {
		return errors.New("invalid label_name")
	}

	return nil
}

func (label *Label) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*label = DefaultLabel()

	type plain Label
	if err := unmarshal((*plain)(label)); err != nil {
		return err
	}

	err := label.Validate()
	if err != nil {
		return err
	}

	return nil
}
