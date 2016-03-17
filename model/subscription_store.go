package model

import (
	"bytes"
)

// S groups all subscription database methods
var S = &subscriptionStore{}

type subscriptionStore struct{}

func (z *subscriptionStore) Delete(tx Transaction, id string) error {
	return delete(tx, entitySubscription, id)
}

// Get returns the subscription with the given compoundID or the given userID and feedID
func (z *subscriptionStore) Get(tx Transaction, id ...string) *Subscription {
	compoundID := ""
	switch len(id) {
	case 1:
		compoundID = id[0]
	case 2:
		compoundID = keyEncode(id...)
	default:
		return nil
	}
	bData := tx.Bucket(bucketData, entitySubscription)
	if data := bData.Get([]byte(compoundID)); data != nil {
		subscription := &Subscription{}
		if err := subscription.decode(data); err == nil {
			return subscription
		}
	}
	return nil
}

func (z *subscriptionStore) GetForUser(tx Transaction, userID string) Subscriptions {
	subscriptions := Subscriptions{}
	min, max := keyMinMax(userID)
	c := tx.Bucket(bucketData, entitySubscription).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		subscription := &Subscription{}
		if err := subscription.decode(v); err == nil {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (z *subscriptionStore) GetForFeed(tx Transaction, feedID string) Subscriptions {
	// index Subscription Feed = FeedID|UserID : UserID|FeedID
	subscriptions := Subscriptions{}
	min, max := keyMinMax(feedID)
	c := tx.Bucket(bucketIndex, entitySubscription, indexSubscriptionFeed).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		compoundID := string(v)
		if subscription := z.Get(tx, compoundID); subscription != nil {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (z *subscriptionStore) New(userID, feedID string) *Subscription {
	return &Subscription{
		UserID: userID,
		FeedID: feedID,
	}
}

func (z *subscriptionStore) Save(tx Transaction, subscription *Subscription) error {
	return save(tx, entitySubscription, subscription)
}
