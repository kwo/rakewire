package fever

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kwo/rakewire/model"
	"github.com/matryer/silk/runner"
	"golang.org/x/net/context"
)

const (
	testUsername = "jeff"
	testPassword = "abcdefg"
)

func TestAuth(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)
	server := newServer(database)
	defer server.Close()

	r := runner.New(t, server.URL)
	r.RunFile("testdata/auth.md")

}

func newServer(database model.Database) *httptest.Server {
	apiFever := New(database)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiFever.ServeHTTPC(context.Background(), w, r)
	})
	return httptest.NewServer(handler)
}

func openTestDatabase(t *testing.T) model.Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	store, err := model.Instance.Open(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	err = store.Update(func(tx model.Transaction) error {
		return populateDatabase(tx)
	})
	if err != nil {
		t.Fatalf("Cannot populate database: %s", err.Error())
	}

	return store

}

func closeTestDatabase(t *testing.T, db model.Database) {

	location := db.Location()

	if err := model.Instance.Close(db); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func getUser(t *testing.T, db model.Database) *model.User {

	var user *model.User

	err := db.Select(func(tx model.Transaction) error {
		user = model.U.GetByUsername(tx, testUsername)
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	return user

}

func populateDatabase(tx model.Transaction) error {

	// add test user
	user := model.U.New(testUsername, testPassword)
	if err := model.U.Save(tx, user); err != nil {
		return err
	}

	// add test groups
	mGroups := []*model.Group{}
	for i := 0; i < 2; i++ {
		g := model.G.New(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := model.G.Save(tx, g); err != nil {
			return err
		}
	}

	// add test feeds
	mFeeds := []*model.Feed{}
	mSubscriptions := []*model.Subscription{}
	for i := 0; i < 4; i++ {
		f := model.F.New(fmt.Sprintf("http://localhost%d", i))
		if err := model.F.Save(tx, f); err != nil {
			return err
		}
		mFeeds = append(mFeeds, f)
		s := model.S.New(user.ID, f.ID)
		s.GroupIDs = append(s.GroupIDs, mGroups[i%2].ID)
		if err := model.S.Save(tx, s); err != nil {
			return err
		}
		mSubscriptions = append(mSubscriptions, s)
	}

	// add test items
	mItems := model.Items{}
	for _, f := range mFeeds {
		now := time.Now().Truncate(time.Second)
		for i := 0; i < 10; i++ {
			item := model.I.New(f.ID, fmt.Sprintf("Item%d", i))
			item.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
			item.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
			mItems = append(mItems, item)
		}
		tr := model.T.New(f.ID)
		tr.StartTime = now
		if err := model.T.Save(tx, tr); err != nil {
			return err
		}
	}
	if err := model.I.SaveAll(tx, mItems); err != nil {
		return err
	}
	if err := model.E.AddItems(tx, mItems); err != nil {
		return err
	}

	// mark entries read
	entries := model.E.Range(tx, user.ID)
	now := time.Now().Truncate(time.Second)
	tRead := now.Add(-6 * 24 * time.Hour).Add(1 * time.Second)
	tStar := now.Add(-8 * 24 * time.Hour).Add(1 * time.Second)
	for _, e := range entries {
		e.Read = e.Updated.Before(tRead)
		e.Star = e.Updated.Before(tStar)
	}
	if err := model.E.SaveAll(tx, entries); err != nil {
		return err
	}

	return nil

}
