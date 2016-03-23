package model

import (
	"encoding/json"
	"time"
)

// GetID returns the unique ID for the object
func (z *Entry) GetID() string {
	return keyEncode(z.UserID, z.ItemID)
}

func (z *Entry) setID(tx Transaction) error {
	return nil
}

func (z *Entry) clear() {
	z.UserID = empty
	z.ItemID = empty
	z.FeedID = empty
	z.Updated = time.Time{}
	z.Read = false
	z.Star = false
}

func (z *Entry) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Entry) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Entry) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexEntryFeedReadUpdated] = []string{z.UserID, z.FeedID, keyEncodeBool(z.Read), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryFeedStarUpdated] = []string{z.UserID, z.FeedID, keyEncodeBool(z.Star), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryFeedUpdated] = []string{z.UserID, z.FeedID, keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryReadUpdated] = []string{z.UserID, keyEncodeBool(z.Read), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryStarUpdated] = []string{z.UserID, keyEncodeBool(z.Star), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryUpdated] = []string{z.UserID, keyEncodeTime(z.Updated), z.ItemID}
	return result
}