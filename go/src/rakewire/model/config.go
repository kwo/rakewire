package model

import (
	"strconv"
)

const (
	configurationBucketName = "Config"
)

// Configuration holds application configuration data
type Configuration struct {
	values map[string]string
}

// NewConfiguration creates a blank configuration
func NewConfiguration() *Configuration {
	return &Configuration{
		values: make(map[string]string),
	}
}

// Load reads configuration values from the database
func (z *Configuration) Load(tx Transaction) error {
	z.values = make(map[string]string)
	b := tx.Bucket(configurationBucketName)
	return b.ForEach(func(k, v []byte) error {
		z.values[string(k)] = string(v)
		return nil
	})
}

// Save saves configuration values to the database
func (z *Configuration) Save(tx Transaction) error {
	b := tx.Bucket(configurationBucketName)
	for key, value := range z.values {
		if err := b.Put([]byte(key), []byte(value)); err != nil {
			return err
		}
	}
	return nil
}

// Get retrieves the named configuration value
func (z *Configuration) Get(key string) string {
	return z.values[key]
}

// GetBool retrieves the named configuration value as a boolean
func (z *Configuration) GetBool(key string) bool {
	if value, ok := z.values[key]; ok {
		return value == "1"
	}
	return false
}

// GetInt retrieves the named configuration value as an integer
func (z *Configuration) GetInt(key string) int {
	if value, ok := z.values[key]; ok {
		result, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return int(result)
		}
	}
	return 0
}

// Set sets a configuration value
func (z *Configuration) Set(key, value string) {
	z.values[key] = value
}

// SetBool sets the named configuration value as a boolean
func (z *Configuration) SetBool(key string, value bool) {
	if value {
		z.values[key] = "1"
	} else {
		z.values[key] = "0"
	}
}

// SetInt sets the named configuration value as an integer
func (z *Configuration) SetInt(key string, value int) {
	z.values[key] = strconv.FormatInt(int64(value), 10)
}
