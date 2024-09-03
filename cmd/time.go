package cmd

import (
	"strconv"
	"strings"
)

func AsSeconds(userTime string) (int, error) {

	multiplier := 1
	value := strings.TrimSpace(userTime)

	if strings.HasSuffix(value, "s") {
		value = userTime[:len(value)-1]
	}

	if strings.HasSuffix(value, "m") {
		multiplier = 60
		value = userTime[:len(value)-1]
	}

	if strings.HasSuffix(value, "h") {
		multiplier = 60 * 60
		value = userTime[:len(value)-1]
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		return -1, err
	}
	return res * multiplier, nil
}

func SecondsToMillis(seconds int) int {
	return seconds * 1000
}

func MinutesToMillis(minutes int) int {
	return minutes * 60 * 1000
}
