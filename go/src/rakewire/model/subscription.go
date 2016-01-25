package model

import (
	"time"
)

//go:generate gokv $GOFILE

// Subscription defines a feed specific to a user.
type Subscription struct {
	ID        uint64
	UserID    uint64 `kv:"+required,Feed:2,User:1"`
	FeedID    uint64 `kv:"+required,Feed:1,User:2"`
	GroupIDs  []uint64
	DateAdded time.Time
	Title     string
	Notes     string
	AutoRead  bool
	AutoStar  bool
	Feed      *Feed `kv:"-"`
}

// NewSubscription associates a feed with a user.
func NewSubscription(userID, feedID uint64) *Subscription {
	return &Subscription{
		UserID:   userID,
		FeedID:   feedID,
		GroupIDs: []uint64{},
	}
}

// AddGroup adds the subscription to the given group.
func (z *Subscription) AddGroup(groupID uint64) {
	if !z.HasGroup(groupID) {
		z.GroupIDs = append(z.GroupIDs, groupID)
	}
}

// RemoveGroup removes the Subscription from the given group.
func (z *Subscription) RemoveGroup(groupID uint64) {
	for i, value := range z.GroupIDs {
		if value == groupID {
			z.GroupIDs = append(z.GroupIDs[:i], z.GroupIDs[i+1:]...)
		}
	}
}

// HasGroup tests if the Subscription belongs to the given group
func (z *Subscription) HasGroup(groupID uint64) bool {
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
