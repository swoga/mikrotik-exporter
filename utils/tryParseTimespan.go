package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var regexTimespan = regexp.MustCompile(`^(?:(\d+)(w))?(?:(\d+)(d))?(?:(\d+)(h))?(?:(\d+)(m))?(?:(\d+)(s))?$`)

func TryParseTimespan(str string) (float64, error) {
	match := regexTimespan.FindStringSubmatch(str)
	if match == nil {
		return 0, errors.New("no regex match")
	}

	var value float64 = 0
	for i := 1; i < len(match); i++ {
		groupValueRaw := match[i]
		i++
		groupSuffix := match[i]

		if groupValueRaw == "" {
			// value is empty if an optional group did not match
			continue
		}
		groupValue, err := strconv.ParseFloat(groupValueRaw, 64)
		if err != nil {
			return 0, err
		}

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
			return 0, fmt.Errorf("invalid timespan suffix: %s", groupSuffix)
		}
	}

	return value, nil
}
