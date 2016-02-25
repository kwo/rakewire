package model

// SubscriptionsByUser retrieves the subscriptions belonging to the user with the Feed populated.
func SubscriptionsByUser(userID string, tx Transaction) (Subscriptions, error) {

	subscriptions := Subscriptions{}

	// subscription index user = UserID|FeedID : SubscriptionID
	min, max := kvKeyMinMax(userID)
	bIndex := tx.Bucket(bucketIndex, subscriptionEntity, subscriptionIndexUser)
	bSubscription := tx.Bucket(bucketData, subscriptionEntity)
	bFeed := tx.Bucket(bucketData, feedEntity)

	err := bIndex.IterateIndex(bSubscription, min, max, func(id string, record Record) error {
		subscription := &Subscription{}
		if err := subscription.deserialize(record); err != nil {
			return err
		}
		if data := bFeed.GetRecord(subscription.FeedID); data != nil {
			feed := &Feed{}
			if err := feed.deserialize(data); err != nil {
				return err
			}
			subscription.Feed = feed
			subscriptions = append(subscriptions, subscription)
		}
		return nil
	})

	return subscriptions, err

}

// SubscriptionsByFeed retrieves the subscriptions associated with the feed.
func SubscriptionsByFeed(feedID string, tx Transaction) (Subscriptions, error) {

	subscriptions := Subscriptions{}

	// subscription index feed = FeedID|UserID : SubscriptionID
	min, max := kvKeyMinMax(feedID)
	bIndex := tx.Bucket(bucketIndex, subscriptionEntity, subscriptionIndexFeed)
	bSubscription := tx.Bucket(bucketData, subscriptionEntity)
	bFeed := tx.Bucket(bucketData, feedEntity)

	err := bIndex.IterateIndex(bSubscription, min, max, func(id string, record Record) error {
		subscription := &Subscription{}
		if err := subscription.deserialize(record); err != nil {
			return err
		}
		if data := bFeed.GetRecord(subscription.FeedID); data != nil {
			feed := &Feed{}
			if err := feed.deserialize(data); err != nil {
				return err
			}
			subscription.Feed = feed
			subscriptions = append(subscriptions, subscription)
		}
		return nil
	})

	return subscriptions, err

}

// Delete removes a subscription from the database.
func (subscription *Subscription) Delete(tx Transaction) error {
	return kvDelete(subscriptionEntity, subscription, tx)
}

// Save saves a user to the database.
func (subscription *Subscription) Save(tx Transaction) error {
	return kvSave(subscriptionEntity, subscription, tx)
}
