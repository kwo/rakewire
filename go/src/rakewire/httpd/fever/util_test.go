package fever

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	"rakewire/logging"
	m "rakewire/model"
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

func makeRequest(user *m.User, target string) ([]byte, error) {

	values := url.Values{}
	if user != nil {
		values.Set(AuthParam, user.FeverHash)
	}
	rsp, err := http.PostForm(target, values)
	if err != nil {
		return nil, err
	} else if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad error code, expected %d, actual %d", http.StatusOK, rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func openDatabase(t *testing.T) (*bolt.Service, string) {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s\n", err.Error())
	}
	testDatabaseFile := f.Name()
	f.Close()

	cfg := db.Configuration{
		Location: testDatabaseFile,
	}
	testDatabase := bolt.NewService(&cfg)
	err = testDatabase.Start()
	if err != nil {
		t.Fatalf("Cannot open database: %s\n", err.Error())
	}

	if err := populateDatabase(testDatabase); err != nil {
		t.Fatalf("Cannot populate database: %s", err.Error())
	}

	return testDatabase, testDatabaseFile

}

func closeDatabase(t *testing.T, database *bolt.Service, testDatabaseFile string) {

	database.Stop()
	if err := os.Remove(testDatabaseFile); err != nil {
		t.Errorf("Cannot delete temp database file: %s", err.Error())
	}

}

func populateDatabase(database *bolt.Service) error {

	// add test user
	user := m.NewUser(testUsername)
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		return err
	}

	// add test groups
	mGroups := []*m.Group{}
	for i := 0; i < 2; i++ {
		g := m.NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := database.GroupSave(g); err != nil {
			return err
		}
	}

	// add test feeds
	mFeeds := []*m.Feed{}
	mUserFeeds := []*m.UserFeed{}
	for i := 0; i < 4; i++ {
		f := m.NewFeed(fmt.Sprintf("http://localhost%d", i))
		if _, err := database.FeedSave(f); err != nil {
			return err
		}
		mFeeds = append(mFeeds, f)
		uf := m.NewUserFeed(user.ID, f.ID)
		uf.GroupIDs = append(uf.GroupIDs, mGroups[i%2].ID)
		database.UserFeedSave(uf)
		mUserFeeds = append(mUserFeeds, uf)
	}

	// add test entries
	for _, f := range mFeeds {
		now := time.Now().Truncate(time.Second)
		for i := 0; i < 10; i++ {
			entry := m.NewEntry(f.ID, fmt.Sprintf("Entry%d", i))
			entry.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
			entry.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
			f.Entries = append(f.Entries, entry)
		}
		if _, err := database.FeedSave(f); err != nil {
			return err
		}
	}

	return nil

}
