package model

import (
	"encoding/json"
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
	Values map[string]string
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

// GetInt64 returns the given value if exists otherwise the default value
func (z *Configuration) GetInt64(name string, defaultValue ...int64) int64 {
	if valueStr, ok := z.Values[name]; ok {
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return value
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

// SetInt64 sets an integer configuration value
func (z *Configuration) SetInt64(name string, value int64) {
	if value != 0 {
		z.Values[name] = strconv.FormatInt(value, 10)
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

// GetID returns the unique ID for the object
func (z *Configuration) GetID() string {
	return idConfig
}

func (z *Configuration) clear() {
	z.Values = make(map[string]string)
}

func (z *Configuration) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Configuration) hasIncrementingID() bool {
	return false
}

func (z *Configuration) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Configuration) indexes() map[string][]string {
	return map[string][]string{}
}

func (z *Configuration) setID(tx Transaction) error {
	return nil
}
