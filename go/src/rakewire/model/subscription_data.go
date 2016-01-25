package model

import (
	"bytes"
	"strconv"
)

// SubscriptionsByUser retrieves the subscriptions belonging to the user with the Feed populated.
func SubscriptionsByUser(userID uint64, tx Transaction) ([]*Subscription, error) {

	var result []*Subscription

	// define index keys
	uf := &Subscription{}
	uf.UserID = userID
	minKeys := uf.IndexKeys()[SubscriptionIndexUser]
	uf.UserID = userID + 1
	nxtKeys := uf.IndexKeys()[SubscriptionIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(SubscriptionEntity).Bucket(SubscriptionIndexUser)
	bSubscription := tx.Bucket(bucketData).Bucket(SubscriptionEntity)
	bFeed := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bSubscription); ok {
			uf := &Subscription{}
			if err := uf.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.Deserialize(data); err != nil {
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
func SubscriptionsByFeed(feedID uint64, tx Transaction) ([]*Subscription, error) {

	var result []*Subscription

	// define index keys
	uf := &Subscription{}
	uf.FeedID = feedID
	minKeys := uf.IndexKeys()[SubscriptionIndexFeed]
	uf.FeedID = feedID + 1
	nxtKeys := uf.IndexKeys()[SubscriptionIndexFeed]

	bIndex := tx.Bucket(bucketIndex).Bucket(SubscriptionEntity).Bucket(SubscriptionIndexFeed)
	bSubscription := tx.Bucket(bucketData).Bucket(SubscriptionEntity)
	bFeed := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bSubscription); ok {
			uf := &Subscription{}
			if err := uf.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.Deserialize(data); err != nil {
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
	return kvDelete(SubscriptionEntity, subscription, tx)
}

// Save saves a user to the database.
func (subscription *Subscription) Save(tx Transaction) error {
	return kvSave(SubscriptionEntity, subscription, tx)
}
