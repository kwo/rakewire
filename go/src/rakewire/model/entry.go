package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Entry defines an item's status for a specific user.
type Entry struct {
	ID             string    `kv:"Read:4,Star:4,User:2"`
	UserID         string    `kv:"+required,Read:1,Star:1,User:1"`
	ItemID         string    `kv:"+required"`
	SubscriptionID string    `kv:"+required"`
	Updated        time.Time `kv:"Read:3,Star:3"`
	IsRead         bool      `kv:"Read:2"`
	IsStar         bool      `kv:"Star:2"`
	Item           *Item     `kv:"-"`
}

func (z *Entry) setIDIfNecessary(fn fnUniqueID) error {
	if z.ID == "0" {
		if _, id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
