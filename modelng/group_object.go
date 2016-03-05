package modelng

import (
	"encoding/json"
)

func (z *Group) getID() string {
	return z.ID
}

func (z *Group) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.Group = config.Sequences.Group + 1
	z.ID = keyEncodeUint(config.Sequences.Group)
	return C.Put(tx, config)
}

func (z *Group) clear() {
	z.ID = empty
	z.UserID = empty
	z.Name = empty
}

func (z *Group) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Group) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Group) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexGroupUserName] = []string{z.UserID, z.Name}
	return result
}
