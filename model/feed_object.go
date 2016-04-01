package model

import (
	"encoding/json"
	"strings"
	"time"
)

// GetID returns the unique ID for the object
func (z *Feed) GetID() string {
	return z.ID
}

func (z *Feed) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.Feed = config.Sequences.Feed + 1
	z.ID = keyEncodeUint(config.Sequences.Feed)
	return C.Put(tx, config)
}

func (z *Feed) clear() {
	z.ID = empty
	z.URL = empty
	z.SiteURL = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.LastUpdated = time.Time{}
	z.NextFetch = time.Time{}
	z.Notes = empty
	z.Title = empty
	z.Status = empty
	z.StatusMessage = empty
	z.StatusSince = time.Time{}
}

func (z *Feed) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Feed) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Feed) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexFeedNextFetch] = []string{keyEncodeTime(z.NextFetch), z.ID}
	result[indexFeedURL] = []string{strings.ToLower(z.URL)}
	return result
}

func (z *Feeds) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Feeds) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}
