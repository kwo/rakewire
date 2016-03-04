package modelng

import (
	"time"
)

const (
	entitySubscription    = "Subscription"
	indexSubscriptionFeed = "Feed"
)

var (
	indexesSubscription = []string{
		indexSubscriptionFeed,
	}
)

// Subscriptions is a collection of Subscription objects.
type Subscriptions []*Subscription

// Subscription defines a feed specific to a user.
type Subscription struct {
	UserID   string    `json:"userId"`
	FeedID   string    `json:"feedId"`
	GroupIDs []string  `json:"groupIds,omitempty"`
	Added    time.Time `json:"added,omitempty"`
	Title    string    `json:"title,omitempty"`
	Notes    string    `json:"notes,omitempty"`
	AutoRead bool      `json:"autoread,omitempty"`
	AutoStar bool      `json:"autostar,omitempty"`
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
