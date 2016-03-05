package modelng

import (
	"testing"
)

func TestSubscriptionSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entitySubscription); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entitySubscription]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestSubscriptions(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	userID := "0000000001"
	feedID := "0000000002"

	// add subscription
	err := db.Update(func(tx Transaction) error {

		subscription := S.New(userID, feedID)
		if err := S.Save(tx, subscription); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error adding subscription: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		subscription := S.Get(tx, userID, feedID)
		if subscription == nil {
			t.Fatal("Nil subscription, expected valid subscription")
		}
		if subscription.UserID != userID {
			t.Errorf("bad userID: %s, expected %s", subscription.UserID, userID)
		}
		if subscription.FeedID != feedID {
			t.Errorf("bad feedID: %s, expected %s", subscription.FeedID, feedID)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting subscription: %s", err.Error())
	}

	// delete subscription
	err = db.Update(func(tx Transaction) error {
		if err := S.Delete(tx, keyEncode(userID, feedID)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting subscription: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		subscription := S.Get(tx, userID, feedID)
		if subscription != nil {
			t.Error("Expected nil subscription")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting subscription: %s", err.Error())
	}

}

func TestSubscriptionsForUserFeed(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	// add subscriptions
	err := db.Update(func(tx Transaction) error {
		for u := 0; u < 3; u++ {
			userID := keyEncodeUint(uint64(u + 1))
			for f := 0; f < 5; f++ {
				feedID := keyEncodeUint(uint64(f + 1))
				subscription := S.New(userID, feedID)
				if err := S.Save(tx, subscription); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error adding subscriptions: %s", err.Error())
	}

	// test by feed
	err = db.Select(func(tx Transaction) error {

		subscriptions := S.GetForFeed(tx, "0000000002")
		if subscriptions == nil {
			t.Fatal("Nil subscriptions, expected valid subscriptions")
		}
		if len(subscriptions) != 3 {
			t.Errorf("bad subscription count %d, expected %d", len(subscriptions), 3)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting subscriptions: %s", err.Error())
	}

	// test by user
	err = db.Select(func(tx Transaction) error {

		subscriptions := S.GetForUser(tx, "0000000002")
		if subscriptions == nil {
			t.Fatal("Nil subscriptions, expected valid subscriptions")
		}
		if len(subscriptions) != 5 {
			t.Errorf("bad subscription count %d, expected %d", len(subscriptions), 5)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting subscriptions: %s", err.Error())
	}

}

func TestSubscriptionGroups(t *testing.T) {

	t.Parallel()

	subscription := S.New(keyEncodeUint(1), keyEncodeUint(1))

	subscription.AddGroup(keyEncodeUint(3))
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.AddGroup(keyEncodeUint(2))
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	subscription.AddGroup(keyEncodeUint(2))
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	if found := subscription.HasGroup(keyEncodeUint(2)); !found {
		t.Error("Group not found: 2")
	}

	if found := subscription.HasGroup(keyEncodeUint(3)); !found {
		t.Error("Group not found: 3")
	}

	if found := subscription.HasGroup(keyEncodeUint(1)); found {
		t.Error("Unexpected group found: 1")
	}

	subscription.RemoveGroup(keyEncodeUint(3))
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.RemoveGroup(keyEncodeUint(2))
	if len(subscription.GroupIDs) != 0 {
		t.Errorf("Bad group count: expcected %d, actual %d", 0, len(subscription.GroupIDs))
	}

}
