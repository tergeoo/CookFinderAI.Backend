package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func ParseIntParam(value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("parameter is required - %s", value)
	}
	return strconv.Atoi(value)
}

func ParseDateParam(value, timeFormat string, defaultDate time.Time) (time.Time, error) {
	if value == "" {
		return defaultDate, nil
	}

	return time.Parse(timeFormat, value)
}

func ParseBoolParam(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean value - %s", value)
	}
	return &boolValue, nil
}

func QueryToMap(r *http.Request) map[string]string {
	query := r.URL.Query()

	params := make(map[string]string)

	for key, values := range query {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params
}
