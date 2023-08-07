package util

import (
	"fmt"
	"strconv"
)

func ParseInt64QueryParam(param string, defaultValue int64, validator func(int64) bool) (int64, error) {
	if param == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid query parameter, not an integer: %v", err)
	}
	if !validator(value) {
		return 0, fmt.Errorf("Invalid query parameter, invalid value: %v", value)
	}
	return value, nil
}
