package modelng

import (
	"fmt"
)

func (z *Config) getID() string {
	return z.Name
}

func (z *Config) setID(fn fnUniqueID) error {
	return nil
}

func (z *Config) encode() ([]byte, error) {
	return keyEncode(z.Name, z.Value), nil
}

func (z *Config) decode(data []byte) error {
	e := keyDecode(data)
	if len(e) != 2 {
		return fmt.Errorf("Invalid blob %s", string(data))
	}
	z.Name = e[0]
	z.Value = e[1]
	return nil
}

func (z *Config) indexes() map[string][]string {
	return map[string][]string{}
}
