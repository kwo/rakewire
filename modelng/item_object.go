package modelng

import (
	"encoding/json"
	"time"
)

func (z *Item) getID() string {
	return z.ID
}

func (z *Item) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.Item = config.Sequences.Item + 1
	z.ID = keyEncodeUint(config.Sequences.Item)
	return C.Put(tx, config)
}

func (z *Item) clear() {
	z.ID = empty
	z.GUID = empty
	z.FeedID = empty
	z.Created = time.Time{}
	z.Updated = time.Time{}
	z.URL = empty
	z.Author = empty
	z.Title = empty
	z.Content = empty
}

func (z *Item) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Item) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Item) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexItemGUID] = []string{z.FeedID, z.GUID}
	return result
}
