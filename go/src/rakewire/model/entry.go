package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Entry defines an item's status for a specific user.
type Entry struct {
	ID             uint64    `kv:"Read:4,Star:4,User:2"`
	UserID         uint64    `kv:"+required,Read:1,Star:1,User:1"`
	ItemID         uint64    `kv:"+required"`
	SubscriptionID uint64    `kv:"+required"`
	Updated        time.Time `kv:"Read:3,Star:3"`
	IsRead         bool      `kv:"Read:2"`
	IsStar         bool      `kv:"Star:2"`
	Item           *Item     `kv:"-"`
}

func (z *Entry) setIDIfNecessary(fn fnNextID, tx Transaction) error {
	if z.ID == 0 {
		if id, err := fn(entryEntity, tx); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
