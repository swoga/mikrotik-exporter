package config

import "github.com/rs/zerolog"

type Labels []Label

func (x Labels) GetByName(name string) (*Label, int) {
	for i, label := range x {
		if label.LabelName == name {
			return &label, i
		}
	}
	return nil, -1
}

func (x *Labels) Add(item Label) {
	tmp := *x
	*x = append(tmp, item)
}

func (x *Labels) RemoveByIndex(i int) {
	tmp := *x
	*x = append(tmp[:i], tmp[i+1:]...)
}

func (x Labels) LabelNames() []string {
	names := make([]string, len(x))
	for i, label := range x {
		names[i] = label.LabelName
	}
	return names
}

func (x Labels) LabelValues(log zerolog.Logger, response map[string]string, variables map[string]string) []string {
	values := make([]string, len(x))
	for i, label := range x {
		labelValue := label.AsString(log, response, variables)
		values[i] = labelValue
	}
	return values
}
