package utils

import (
	"regexp"
	"strconv"
)

var regexAZ = regexp.MustCompile("[a-zA-Z]+")

func TryParseDouble(valueP *string) (float64, bool) {
	if valueP == nil {
		return 0, false
	}
	value := *valueP

	// remove potential units
	value = regexAZ.ReplaceAllLiteralString(value, "")
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}
