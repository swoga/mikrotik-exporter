package utils

import (
	"strings"

	"github.com/rs/zerolog"
)

func Substitute(log zerolog.Logger, value string, variables map[string]string) string {
	newValue := value

	for varName, varValue := range variables {
		newValue = strings.ReplaceAll(newValue, "{"+varName+"}", varValue)
	}

	if value != newValue {
		log.Trace().Str("input", value).Str("output", newValue).Msg("substitute value")
	}

	return newValue
}
