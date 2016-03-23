package model

import (
	"encoding/json"
)

// GetID returns the unique ID for the object
func (z *Config) GetID() string {
	return idConfig
}

func (z *Config) setID(tx Transaction) error {
	return nil
}

func (z *Config) clear() {
	z.ID = empty
	z.LoggingLevel = empty
	z.Sequences = sequences{}
}

func (z *Config) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Config) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Config) indexes() map[string][]string {
	return map[string][]string{}
}
