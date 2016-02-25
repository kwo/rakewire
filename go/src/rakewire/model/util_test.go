package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"rakewire/logging"
	"testing"
	"time"
)

const (
	testUsername = "jeff"
)

func TestMain(m *testing.M) {
	cfg := &logging.Configuration{Level: logging.LogWarn}
	cfg.Init()
	status := m.Run()
	os.Exit(status)
}

func openTestDatabase(t *testing.T, flags ...bool) Database {

	flagPopulateDatabase := len(flags) > 0 && flags[0]

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	boltDB, err := OpenDatabase(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	if flagPopulateDatabase {
		err = boltDB.Update(func(tx Transaction) error {
			return populateDatabase(t, tx)
		})
		if err != nil {
			t.Fatalf("Cannot populate database: %s", err.Error())
		}
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d Database) {

	location := d.Location()

	if err := CloseDatabase(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func populateDatabase(t *testing.T, tx Transaction) error {

	// add test user
	user := NewUser(testUsername)
	user.SetPassword("abcdefg")
	if err := user.Save(tx); err != nil {
		return err
	}

	// add test groups
	mGroups := []*Group{}
	for i := 0; i < 2; i++ {
		g := NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := g.Save(tx); err != nil {
			return err
		}
	}

	// add test feeds and subscriptions
	mFeeds := []*Feed{}
	mSubscriptions := []*Subscription{}
	for i := 0; i < 4; i++ {
		f := NewFeed(fmt.Sprintf("http://localhost%d", i))
		if _, err := f.Save(tx); err != nil {
			return err
		}
		mFeeds = append(mFeeds, f)
		uf := NewSubscription(user.ID, f.ID)
		uf.GroupIDs = append(uf.GroupIDs, mGroups[i%2].ID)
		if err := uf.Save(tx); err != nil {
			return err
		}
		mSubscriptions = append(mSubscriptions, uf)
	}

	// add test items
	for n, f := range mFeeds {
		now := time.Now().Truncate(time.Second)
		for i := 0; i < 10; i++ {
			item := NewItem(f.ID, fmt.Sprintf("Feed%dItem%d", n, i))
			item.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
			item.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
			//t.Logf("Item: %v", item)
			f.Items = append(f.Items, item)
		}
		f.Transmission = NewTransmission(f.ID)
		f.Transmission.StartTime = now
		if _, err := f.Save(tx); err != nil {
			return err
		}
	}

	// mark entries read
	entries, err := EntriesGetAll(user.ID, tx)
	if err != nil {
		return err
	}
	now := time.Now().Truncate(time.Second)
	tRead := now.Add(-6 * 24 * time.Hour).Add(1 * time.Second)
	tStar := now.Add(-8 * 24 * time.Hour).Add(1 * time.Second)
	for _, ue := range entries {
		ue.IsRead = ue.Updated.Before(tRead)
		ue.IsStar = ue.Updated.Before(tStar)
	}
	if err := EntriesSave(entries, tx); err != nil {
		return err
	}

	return nil

}
