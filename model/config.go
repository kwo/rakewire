package model

import (
	"strconv"
)

const (
	entityConfig = "Config"
	idConfig     = "configuration"
)

var (
	indexesConfig = []string{}
)

// Configuration defines the application configurtion.
type Configuration struct {
	ID        string
	Sequences sequences
	Values    map[string]string
}

type sequences struct {
	User         uint64
	Feed         uint64
	Item         uint64
	Group        uint64
	Transmission uint64
}

// GetBool returns the given value if exists otherwise the default value
func (z *Configuration) GetBool(name string, defaultValue ...bool) bool {
	if valueStr, ok := z.Values[name]; ok {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// SetBool sets a boolean configuration value
func (z *Configuration) SetBool(name string, value bool) {
	if value {
		z.Values[name] = strconv.FormatBool(value)
	} else {
		delete(z.Values, name)
	}
}

// GetInt returns the given value if exists otherwise the default value
func (z *Configuration) GetInt(name string, defaultValue ...int) int {
	if valueStr, ok := z.Values[name]; ok {
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return int(value)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// SetInt sets an integer configuration value
func (z *Configuration) SetInt(name string, value int) {
	if value != 0 {
		z.Values[name] = strconv.FormatInt(int64(value), 10)
	} else {
		delete(z.Values, name)
	}
}

// GetStr returns the given value if exists otherwise the default value
func (z *Configuration) GetStr(name string, defaultValue ...string) string {
	if value, ok := z.Values[name]; ok {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return empty
}

// SetStr sets an integer configuration value
func (z *Configuration) SetStr(name string, value string) {
	if value != empty {
		z.Values[name] = value
	} else {
		delete(z.Values, name)
	}
}
