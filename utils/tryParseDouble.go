package utils

import (
	"errors"
	"regexp"
	"strconv"
)

var regexAZ = regexp.MustCompile("[a-zA-Z]+")

func TryParseDouble(valueP *string) (float64, error) {
	if valueP == nil {
		return 0, errors.New("input nil")
	}
	value := *valueP

	// remove potential units
	value = regexAZ.ReplaceAllLiteralString(value, "")
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
