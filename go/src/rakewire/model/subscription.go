package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Subscription defines a feed specific to a user.
type Subscription struct {
	ID        string
	UserID    string `kv:"+required,Feed:2,User:1"`
	FeedID    string `kv:"+required,Feed:1,User:2"`
	GroupIDs  []string
	DateAdded time.Time
	Title     string
	Notes     string
	AutoRead  bool
	AutoStar  bool
	Feed      *Feed `kv:"-"`
}

// NewSubscription associates a feed with a user.
func NewSubscription(userID, feedID string) *Subscription {
	return &Subscription{
		UserID:   userID,
		FeedID:   feedID,
		GroupIDs: []string{},
	}
}

func (z *Subscription) setIDIfNecessary(fn fnUniqueID) error {
	if z.ID == "0" {
		if _, id, err := fn(); err == nil {
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
