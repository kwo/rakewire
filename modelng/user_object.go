package modelng

import (
	"encoding/json"
	"strings"
)

func (z *User) getID() string {
	return z.ID
}

func (z *User) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.User = config.Sequences.User + 1
	z.ID = keyEncodeUint(config.Sequences.User)
	return C.Put(config, tx)
}

func (z *User) clear() {
	z.ID = empty
	z.Username = empty
	z.PasswordHash = empty
	z.FeverHash = empty
}

func (z *User) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *User) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *User) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexUserUsername] = []string{strings.ToLower(z.Username)}
	result[indexUserFeverhash] = []string{z.FeverHash}
	return result
}
