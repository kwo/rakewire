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
	"rakewire/logging"
	"rakewire/middleware"
	"rakewire/model"
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

func newServer(database model.Database) *httptest.Server {
	apiFever := NewAPI("/fever", database)
	return httptest.NewServer(middleware.Adapt(apiFever.Router(), middleware.NoCache(), gorillaHandlers.CompressHandler))
}

func makeRequest(user *model.User, target string, formValues ...string) ([]byte, error) {

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

func openTestDatabase(t *testing.T) model.Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	boltDB, err := model.OpenDatabase(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	err = boltDB.Update(func(tx model.Transaction) error {
		return populateDatabase(tx)
	})
	if err != nil {
		t.Fatalf("Cannot populate database: %s", err.Error())
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d model.Database) {

	location := d.Location()

	if err := model.CloseDatabase(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func populateDatabase(tx model.Transaction) error {

	// add test user
	user := model.NewUser(testUsername)
	user.SetPassword("abcdefg")
	if err := user.Save(tx); err != nil {
		return err
	}

	// add test groups
	mGroups := []*model.Group{}
	for i := 0; i < 2; i++ {
		g := model.NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := g.Save(tx); err != nil {
			return err
		}
	}

	// add test feeds
	mFeeds := []*model.Feed{}
	mSubscriptions := []*model.Subscription{}
	for i := 0; i < 4; i++ {
		f := model.NewFeed(fmt.Sprintf("http://localhost%d", i))
		if _, err := f.Save(tx); err != nil {
			return err
		}
		mFeeds = append(mFeeds, f)
		uf := model.NewSubscription(user.ID, f.ID)
		uf.GroupIDs = append(uf.GroupIDs, mGroups[i%2].ID)
		if err := uf.Save(tx); err != nil {
			return err
		}
		mSubscriptions = append(mSubscriptions, uf)
	}

	// add test items
	for _, f := range mFeeds {
		now := time.Now().Truncate(time.Second)
		for i := 0; i < 10; i++ {
			item := model.NewItem(f.ID, fmt.Sprintf("Item%d", i))
			item.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
			item.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
			f.Items = append(f.Items, item)
		}
		f.Transmission = model.NewTransmission(f.ID)
		f.Transmission.StartTime = now
		if _, err := f.Save(tx); err != nil {
			return err
		}
	}

	// mark entries read
	entries, err := model.EntriesGetNext(user.ID, 0, 0, tx)
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
	if err := model.EntriesSave(entries, tx); err != nil {
		return err
	}

	return nil

}
