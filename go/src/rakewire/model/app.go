package model

import (
	"fmt"
	"strconv"
	"time"
)

const (
	empty      = ""
	timeFormat = time.RFC3339Nano
)

// application level variables
var (
	BuildHash string
	BuildTime string
	Version   string
)

func getBool(fieldName string, values map[string]string, errors []error) bool {
	result, err := strconv.ParseBool(values[fieldName])
	if err != nil {
		errors = append(errors, err)
		return false
	}
	return result
}

func getDuration(fieldName string, values map[string]string, errors []error) time.Duration {
	var result time.Duration
	if value, ok := values[fieldName]; ok {
		t, err := time.ParseDuration(value)
		if err != nil {
			errors = append(errors, err)
		} else {
			result = t
		}
	}
	return result
}

func getInt(fieldName string, values map[string]string, errors []error) int {
	result, err := strconv.ParseInt(values[fieldName], 10, 64)
	if err != nil {
		errors = append(errors, err)
		return 0
	}
	return int(result)
}

func getString(fieldName string, values map[string]string, errors []error) string {
	return values[fieldName]
}

func getTime(fieldName string, values map[string]string, errors []error) time.Time {
	result := time.Time{}
	if value, ok := values[fieldName]; ok {
		t, err := time.Parse(timeFormat, value)
		if err != nil {
			errors = append(errors, err)
		} else {
			result = t
		}
	}
	return result
}

func getUint(fieldName string, values map[string]string, errors []error) uint64 {
	result, err := strconv.ParseUint(values[fieldName], 10, 64)
	if err != nil {
		errors = append(errors, err)
		return 0
	}
	return result
}

func setBool(value bool, fieldName string, result map[string]string) {
	if value {
		result[fieldName] = strconv.FormatBool(value)
	}
}

func setDuration(value time.Duration, fieldName string, result map[string]string) {
	if value != 0 {
		result[fieldName] = value.String()
	}
}

func setInt(value int, fieldName string, result map[string]string) {
	if value != 0 {
		result[fieldName] = strconv.FormatInt(int64(value), 10)
	}
}

func setString(value string, fieldName string, result map[string]string) {
	if value != empty {
		result[fieldName] = value
	}
}

func setTime(value time.Time, fieldName string, result map[string]string) {
	if !value.IsZero() {
		result[fieldName] = value.Format(timeFormat)
	}
}

func setUint(value uint64, fieldName string, result map[string]string) {
	if value != 0 {
		result[fieldName] = fmt.Sprintf("%05d", value)
	}
}
