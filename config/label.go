package config

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Label struct {
	Param Param `yaml:",inline"`

	LabelName string `yaml:"label_name"`
}

func DefaultLabel() Label {
	return Label{}
}

var regexValidMetricAndLabelName = regexp.MustCompile("^[a-zA-Z_:][a-zA-Z0-9_:]*$")

func (label *Label) GetName() string {
	if label.LabelName != "" {
		return label.LabelName
	}
	return strings.ReplaceAll(label.Param.ParamName, "-", "_")
}

func (label *Label) Validate() error {
	if label.GetName() == "" {
		return errors.New("require param_name or label_name")
	}
	if !regexValidMetricAndLabelName.MatchString(label.GetName()) {
		return errors.New("invalid label_name")
	}

	return nil
}

func (label *Label) UnmarshalYAML(node *yaml.Node) error {
	*label = DefaultLabel()

	type plain Label
	if err := node.Decode((*plain)(label)); err != nil {
		return err
	}

	err := label.Validate()
	if err != nil {
		return err
	}

	return nil
}
