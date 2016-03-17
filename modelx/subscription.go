package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Subscription defines a feed specific to a user.
type Subscription struct {
	ID       string    `json:"id"`
	UserID   string    `json:"userID" kv:"+required,Feed:2,User:1"`
	FeedID   string    `json:"feedID" kv:"+required,Feed:1,User:2"`
	GroupIDs []string  `json:"groupIDs,omitempty"`
	Added    time.Time `json:"added,omitempty"`
	Title    string    `json:"title,omitempty"`
	Notes    string    `json:"notes,omitempty"`
	AutoRead bool      `json:"autoread,omitempty"`
	AutoStar bool      `json:"autostar,omitempty"`
	Feed     *Feed     `json:"-" kv:"-"`
}

// NewSubscription associates a feed with a user.
func NewSubscription(userID, feedID string) *Subscription {
	return &Subscription{
		UserID:   userID,
		FeedID:   feedID,
		GroupIDs: []string{},
	}
}

func (z *Subscription) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}

// AddGroup adds the subscription to the given group.
func (z *Subscription) AddGroup(groupID string) {
	if !z.HasGroup(groupID) {
		z.GroupIDs = append(z.GroupIDs, groupID)
	}
}

// RemoveGroup removes the Subscription from the given group.
func (z *Subscription) RemoveGroup(groupID string) {
	for i, value := range z.GroupIDs {
		if value == groupID {
			z.GroupIDs = append(z.GroupIDs[:i], z.GroupIDs[i+1:]...)
		}
	}
}

// HasGroup tests if the Subscription belongs to the given group
func (z *Subscription) HasGroup(groupID string) bool {
	result := false
	if len(z.GroupIDs) > 0 {
		for _, value := range z.GroupIDs {
			if value == groupID {
				return true
			}
		}
	}
	return result
}
