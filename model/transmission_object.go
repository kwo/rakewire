package model

import (
	"encoding/json"
	"time"
)

func (z *Transmission) GetID() string {
	return z.ID
}

func (z *Transmission) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.Transmission = config.Sequences.Transmission + 1
	z.ID = keyEncodeUint(config.Sequences.Transmission)
	return C.Put(tx, config)
}

func (z *Transmission) clear() {
	z.ID = empty
	z.FeedID = empty
	z.Duration = 0
	z.Result = empty
	z.ResultMessage = empty
	z.StartTime = time.Time{}
	z.URL = empty
	z.ContentLength = 0
	z.ContentType = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.StatusCode = 0
	z.UsesGzip = false
	z.Flavor = empty
	z.Generator = empty
	z.Title = empty
	z.LastUpdated = time.Time{}
	z.ItemCount = 0
	z.NewItems = 0
}

func (z *Transmission) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Transmission) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Transmission) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexTransmissionTime] = []string{keyEncodeTime(z.StartTime), z.ID}
	result[indexTransmissionFeedTime] = []string{z.FeedID, keyEncodeTime(z.StartTime)}
	return result
}
