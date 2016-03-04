package modelng

import (
	"bytes"
)

// S groups all subscription database methods
var S = &subscriptionStore{}

type subscriptionStore struct{}

func (z *subscriptionStore) Delete(id string, tx Transaction) error {
	return delete(entitySubscription, id, tx)
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

func (z *subscriptionStore) GetByIDs(userID, feedID string, tx Transaction) *Subscription {
	return S.Get(keyEncode(userID, feedID), tx)
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

func (z *subscriptionStore) GetForFeed(feedID string, tx Transaction) Subscriptions {
	// index Subscription Feed = FeedID|UserID : UserID|FeedID
	subscriptions := Subscriptions{}
	min, max := keyMinMax(feedID)
	c := tx.Bucket(bucketIndex, entitySubscription, indexSubscriptionFeed).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		compoundID := string(v)
		if subscription := S.Get(compoundID, tx); subscription != nil {
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

func (z *subscriptionStore) Save(subscription *Subscription, tx Transaction) error {
	return save(entitySubscription, subscription, tx)
}
