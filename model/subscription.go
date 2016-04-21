package model

import (
	"encoding/json"
	"sort"
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

// GetID returns the unique ID for the object
func (z *Subscription) GetID() string {
	return keyEncode(z.UserID, z.FeedID)
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

// RemoveGroup removes the Subscription from the given group.
func (z *Subscription) RemoveGroup(groupID string) {
	for i, value := range z.GroupIDs {
		if value == groupID {
			z.GroupIDs = append(z.GroupIDs[:i], z.GroupIDs[i+1:]...)
		}
	}
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

func (z *Subscription) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Subscription) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Subscription) hasIncrementingID() bool {
	return false
}

func (z *Subscription) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexSubscriptionFeed] = []string{z.FeedID, z.UserID}
	return result
}

func (z *Subscription) setID(tx Transaction) error {
	return nil
}

// Subscriptions is a collection of Subscription objects.
type Subscriptions []*Subscription

func (z Subscriptions) Len() int      { return len(z) }
func (z Subscriptions) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Subscriptions) Less(i, j int) bool {
	return z[i].Added.Before(z[j].Added)
}

// ByFeedID groups elements in the Subscriptions collection by FeedID
func (z Subscriptions) ByFeedID() map[string]Subscriptions {
	result := make(map[string]Subscriptions)
	for _, subscription := range z {
		subscriptions := result[subscription.FeedID]
		subscriptions = append(subscriptions, subscription)
		result[subscription.FeedID] = subscriptions
	}
	return result
}

// ByTitle groups elements in the Subscriptions collection by Name
func (z Subscriptions) ByTitle() map[string]*Subscription {
	result := make(map[string]*Subscription)
	for _, subscription := range z {
		result[subscription.Title] = subscription
	}
	return result
}

// SortByAddedDate sort collection by AddedDate
func (z Subscriptions) SortByAddedDate() {
	sort.Stable(z)
}

// WithGroup creates a new Subscriptions collection containing only subscriptions with the given groupID.
func (z Subscriptions) WithGroup(groupID string) Subscriptions {
	result := Subscriptions{}
	for _, subscription := range z {
		if subscription.HasGroup(groupID) {
			result = append(result, subscription)
		}
	}
	return result
}

func (z *Subscriptions) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Subscriptions) encode() ([]byte, error) {
	return json.Marshal(z)
}
