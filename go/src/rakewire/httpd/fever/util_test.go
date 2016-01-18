package fever

import (
	"compress/gzip"
	"fmt"
	gorillaHandlers "github.com/gorilla/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	"rakewire/logging"
	"rakewire/middleware"
	"rakewire/model"
	m "rakewire/model"
	"strings"
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

func newServer(database *bolt.Service) *httptest.Server {
	apiFever := NewAPI("/fever", database)
	return httptest.NewServer(middleware.Adapt(apiFever.Router(), middleware.NoCache(), gorillaHandlers.CompressHandler))
}

func makeRequest(user *m.User, target string, formValues ...string) ([]byte, error) {

	values := url.Values{}
	if user != nil {
		values.Set(AuthParam, user.FeverHash)
	}

	if len(formValues) != 0 {
		if len(formValues)%2 != 0 {
			return nil, fmt.Errorf("form values must be pairs, %d elements found", len(formValues))
		}
		for i, formValue := range formValues {
			if i%2 != 0 {
				continue
			}
			values.Set(formValue, formValues[i+1])
		}
	}

	client := http.Client{}
	req, err := http.NewRequest(mPost, target, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set(hContentType, "application/x-www-form-urlencoded")
	req.Header.Set(hAcceptEncoding, "gzip")

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad error code, expected %d, actual %d", http.StatusOK, rsp.StatusCode)
	} else if rsp.Header.Get(hContentType) != mimeJSON {
		return nil, fmt.Errorf("Bad content type, expected %d, actual %d", mimeJSON, rsp.Header.Get(hContentType))
	} else if rsp.Header.Get(hContentEncoding) != "gzip" {
		return nil, fmt.Errorf("Bad content encoding, expected %s, actual %s", "gzip", rsp.Header.Get(hContentEncoding))
	}

	gzipReader, err := gzip.NewReader(rsp.Body)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(gzipReader)
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

	err := database.Update(func(tx model.Transaction) error {
		return user.Save(tx)
	})
	if err != nil {
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
		f.Attempt = m.NewFeedLog(f.ID)
		f.Attempt.StartTime = now
		if _, err := database.FeedSave(f); err != nil {
			return err
		}
	}

	// mark entries read
	userEntries, err := database.UserEntryGetNext(user.ID, 0, 0)
	if err != nil {
		return err
	}
	now := time.Now().Truncate(time.Second)
	tRead := now.Add(-6 * 24 * time.Hour).Add(1 * time.Second)
	tStar := now.Add(-8 * 24 * time.Hour).Add(1 * time.Second)
	for _, ue := range userEntries {
		ue.IsRead = ue.Updated.Before(tRead)
		ue.IsStar = ue.Updated.Before(tStar)
	}
	if err := database.UserEntrySave(userEntries); err != nil {
		return err
	}

	return nil

}
