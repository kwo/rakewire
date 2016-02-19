package model

import (
	"bytes"
)

// SubscriptionsByUser retrieves the subscriptions belonging to the user with the Feed populated.
func SubscriptionsByUser(userID string, tx Transaction) (Subscriptions, error) {

	result := Subscriptions{}

	// subscription index user = UserID|FeedID : SubscriptionID
	min, max := kvKeyMinMax(userID)
	bIndex := tx.Bucket(bucketIndex).Bucket(subscriptionEntity).Bucket(subscriptionIndexUser)
	bSubscription := tx.Bucket(bucketData).Bucket(subscriptionEntity)
	bFeed := tx.Bucket(bucketData).Bucket(feedEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		subscriptionID := string(v)

		if data, ok := kvGet(subscriptionID, bSubscription); ok {
			uf := &Subscription{}
			if err := uf.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.deserialize(data); err != nil {
					return nil, err
				}
				uf.Feed = f
				result = append(result, uf)
			}
		}

	}

	return result, nil

}

// SubscriptionsByFeed retrieves the subscriptions associated with the feed.
func SubscriptionsByFeed(feedID string, tx Transaction) (Subscriptions, error) {

	result := Subscriptions{}

	// subscription index feed = FeedID|UserID : SubscriptionID
	min, max := kvKeyMinMax(feedID)
	bIndex := tx.Bucket(bucketIndex).Bucket(subscriptionEntity).Bucket(subscriptionIndexFeed)
	bSubscription := tx.Bucket(bucketData).Bucket(subscriptionEntity)
	bFeed := tx.Bucket(bucketData).Bucket(feedEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		subscriptionID := string(v)

		if data, ok := kvGet(subscriptionID, bSubscription); ok {
			uf := &Subscription{}
			if err := uf.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.deserialize(data); err != nil {
					return nil, err
				}
				uf.Feed = f
				result = append(result, uf)
			}
		}

	}

	return result, nil

}

// Delete removes a subscription from the database.
func (subscription *Subscription) Delete(tx Transaction) error {
	return kvDelete(subscriptionEntity, subscription, tx)
}

// Save saves a user to the database.
func (subscription *Subscription) Save(tx Transaction) error {
	return kvSave(subscriptionEntity, subscription, tx)
}
