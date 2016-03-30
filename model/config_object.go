package model

import (
	"encoding/json"
)

// GetID returns the unique ID for the object
func (z *Configuration) GetID() string {
	return idConfig
}

func (z *Configuration) setID(tx Transaction) error {
	return nil
}

func (z *Configuration) clear() {
	z.Sequences = sequences{}
	z.Values = make(map[string]string)
}

func (z *Configuration) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Configuration) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Configuration) indexes() map[string][]string {
	return map[string][]string{}
}
