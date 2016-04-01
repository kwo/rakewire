package model

import (
	"encoding/json"
	"time"
)

func (z *Subscriptions) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Subscriptions) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

// GetID returns the unique ID for the object
func (z *Subscription) GetID() string {
	return keyEncode(z.UserID, z.FeedID)
}

func (z *Subscription) setID(tx Transaction) error {
	return nil
}

func (z *Subscription) clear() {
	z.UserID = empty
	z.FeedID = empty
	z.GroupIDs = []string{}
	z.Added = time.Time{}
	z.Title = empty
	z.Notes = empty
	z.AutoRead = false
	z.AutoStar = false
}

func (z *Subscription) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Subscription) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Subscription) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexSubscriptionFeed] = []string{z.FeedID, z.UserID}
	return result
}
