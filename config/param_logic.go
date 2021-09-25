package config

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/swoga/mikrotik-exporter/utils"
)

func (param *Param) PreprocessValue(log zerolog.Logger, response map[string]string, variables map[string]string) (string, bool) {
	var value string

	if param.Value != "" {
		log.Trace().Str("value", param.Value).Msg("static parameter")
		value = utils.Substitute(log, param.Value, variables)
	} else {
		apiWord, isInResponse := response[param.ParamName]

		if !isInResponse {
			log.Trace().Msg("field not found in API response")
			return "", false
		}
		value = apiWord
		log.Trace().Str("name", param.ParamName).Str("value", value).Msg("got word from response")
	}

	remappedValue, hasStaticRemap := param.RemapValues[value]
	if hasStaticRemap {
		if remappedValue == nil {
			log.Trace().Msg("remapped to null")
			return "", false
		}

		log.Trace().Str("value", value).Msg("remapped")
		return *remappedValue, true
	}

	for regex, remappedValue := range param.remapValueReParsed {
		if !regex.MatchString(value) {
			continue
		}
		if remappedValue == nil {
			log.Trace().Msg("regex remapped to null")
			return "", false
		}
		value = regex.ReplaceAllString(value, *remappedValue)

		log.Trace().Str("value", value).Msg("regex remapped")
		return value, true
	}

	return value, true
}

var regexTimespan = regexp.MustCompile(`(?:(\d+)(w))?(?:(\d+)(d))?(?:(\d+)(h))?(?:(\d+)(m))?(?:(\d+)(s))?`)

func (param *Param) tryGetValue(log zerolog.Logger, response map[string]string, variables map[string]string) (float64, bool) {
	word, isStaticOrInResponse := param.PreprocessValue(log, response, variables)

	// no static value or not found in API response
	if !isStaticOrInResponse {
		log.Trace().Msg("neither static nor in response")

		if param.Default == "" {
			log.Trace().Msg("no default value set")
			return 0, false
		}

		log.Trace().Msg("use default")
		word = utils.Substitute(log, param.Default, variables)
		value, ok := utils.TryParseDouble(&word)
		if !ok {
			log.Error().Str("word", word).Msg("failed to parse default to float")
		}

		return value, true
	}

	parseLog := log.With().Str("word", word).Logger()

	parseLog.Trace().Str("param_type", param.ParamType).Msg("parse as")

	switch param.ParamType {
	case PARAM_TYPE_INT:
		value, ok := utils.TryParseDouble(&word)
		if !ok {
			parseLog.Error().Msg("failed to parse value to float")
			return 0, false
		}
		return value, true
	case PARAM_TYPE_BOOL:
		wordLowers := strings.ToLower(word)
		value := wordLowers == "true" || wordLowers == "yes"

		if param.Negate {
			value = !value
		}

		if value {
			return 1, true
		} else {
			return 0, true
		}
	case PARAM_TYPE_TIMESPAN:
		match := regexTimespan.FindStringSubmatch(word)
		if match == nil {
			parseLog.Error().Msg("failed to parse timespan")
			return 0, false
		}

		var value float64 = 0
		for i := 1; i < len(match); i++ {
			groupValue, err := strconv.ParseFloat(match[i], 64)
			if err != nil {
				parseLog.Err(err).Str("value", match[i]).Msg("failed to parse timespan value")
				return 0, false
			}
			i++
			groupSuffix := match[i]

			switch groupSuffix {
			case "w":
				value += groupValue * 604800
			case "d":
				value += groupValue * 86400
			case "h":
				value += groupValue * 3600
			case "m":
				value += groupValue * 60
			case "s":
				value += groupValue
			default:
				parseLog.Panic().Str("suffix", groupSuffix).Msg("invalid timespan suffix")
			}
		}

		return value, true
	case PARAM_TYPE_DATETIME:
		dateTime, err := time.Parse("Jan/02/2006 15:04:05", word)
		if err != nil {
			parseLog.Err(err).Msg("failed to parse datetime")
			return 0, false
		}

		switch param.DateTimeType {
		case PARAM_DATETYPE_FROM_NOW:
			return time.Since(dateTime).Seconds(), true
		case PARAM_DATETYPE_TO_NOW:
			return time.Until(dateTime).Seconds(), true
		case PARAM_DATETYPE_TIMESTAMP:
			return float64(dateTime.Unix()), true
		default:
			parseLog.Panic().Str("datetime_type", param.DateTimeType).Msg("invalid datetime_type")
			return 0, false
		}
	default:
		parseLog.Panic().Str("ParamType", param.ParamType).Msg("invalid ParamType")
		return 0, false
	}
}
