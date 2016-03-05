package modelng

import (
	"bytes"
)

// S groups all subscription database methods
var S = &subscriptionStore{}

type subscriptionStore struct{}

func (z *subscriptionStore) Delete(tx Transaction, id string) error {
	return delete(tx, entitySubscription, id)
}

func (z *subscriptionStore) Get(id string, tx Transaction) *Subscription {
	bData := tx.Bucket(bucketData, entitySubscription)
	if data := bData.Get([]byte(id)); data != nil {
		subscription := &Subscription{}
		if err := subscription.decode(data); err == nil {
			return subscription
		}
	}
	return nil
}

func (z *subscriptionStore) GetByIDs(tx Transaction, userID, feedID string) *Subscription {
	return z.Get(keyEncode(userID, feedID), tx)
}

func (z *subscriptionStore) GetForUser(userID string, tx Transaction) Subscriptions {
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
		if subscription := z.Get(compoundID, tx); subscription != nil {
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
