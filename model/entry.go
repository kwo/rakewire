package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Entry defines an item's status for a specific user.
type Entry struct {
	ID             string    `json:"id" kv:"Read:4,Star:4,User:2"`
	UserID         string    `json:"userID" kv:"+required,Read:1,Star:1,User:1"`
	ItemID         string    `json:"itemID" kv:"+required"`
	SubscriptionID string    `json:"subscriptionID" kv:"+required"` // TODO: necessary?
	Updated        time.Time `json:"updated" kv:"Read:3,Star:3"`
	IsRead         bool      `json:"read" kv:"Read:2"`
	IsStar         bool      `json:"star" kv:"Star:2"`
	Item           *Item     `json:"-" kv:"-"`
}

// NewEntry returns a new Entry object
func NewEntry(userID, itemID, subscriptionID string) *Entry {

	return &Entry{
		UserID:         userID,
		ItemID:         itemID,
		SubscriptionID: subscriptionID,
	}

}

func (z *Entry) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
