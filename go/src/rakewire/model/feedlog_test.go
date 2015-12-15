package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewFeedLog(t *testing.T) {

	t.Parallel()

	fl := NewFeedLog("123")

	if fl == nil {
		t.Fatalf("FeedLog factory not returning a valid feedlog")
	}

	if fl.ID != 0 {
		t.Errorf("Factory method should not set th ID, expected %d, actual %d", 0, fl.ID)
	}

	if fl.FeedID != "123" {
		t.Errorf("Factory method not setting FeedID properly, expected %s, actual %s", "123", fl.FeedID)
	}

}

func TestFeedLogSerialize(t *testing.T) {

	t.Parallel()

	fl := getNewFeedLog()
	validateFeedLog(t, fl)

	data := fl.Serialize()
	if data == nil {
		t.Fatal("FeedLog serialize returned a nil map")
	}

	fl2 := &FeedLog{}
	if err := fl2.Deserialize(data); err != nil {
		t.Fatalf("FeedLog deserialize returned an error: %s", err.Error())
	}
	validateFeedLog(t, fl2)

}

func TestFeedLogJSON(t *testing.T) {

	t.Parallel()

	fl := getNewFeedLog()
	validateFeedLog(t, fl)

	data, err := json.Marshal(fl)
	assertNoError(t, err)
	assertNotNil(t, data)

	fl2 := &FeedLog{}
	err = json.Unmarshal(data, fl2)
	assertNoError(t, err)
	validateFeedLog(t, fl)

}

func getNewFeedLog() *FeedLog {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	fl := NewFeedLog("123")
	fl.ContentLength = 0
	fl.ContentType = "text/plain"
	fl.Duration = 6 * time.Second
	fl.ETag = "etag"
	fl.Flavor = "flavor"
	fl.Generator = ""
	//fl.LastModified = dt
	fl.LastUpdated = dt
	fl.Result = FetchResultOK
	fl.ResultMessage = "message"
	fl.StartTime = dt
	fl.StatusCode = 200
	fl.Title = "title"
	fl.URL = "url"
	fl.UsesGzip = false
	fl.NewEntries = 2

	return fl

}

func validateFeedLog(t *testing.T, fl *FeedLog) {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	assertNotNil(t, fl)

	if fl.ID != 0 {
		t.Errorf("FeedLog ID incorrect, expected %d, actual %d", 0, fl.ID)
	}

	assertEqual(t, 0, fl.ContentLength)
	assertEqual(t, "text/plain", fl.ContentType)
	assertEqual(t, 6*time.Second, fl.Duration)
	assertEqual(t, "etag", fl.ETag)
	assertEqual(t, "123", fl.FeedID)
	assertEqual(t, "flavor", fl.Flavor)
	assertEqual(t, "", fl.Generator)
	assertEqual(t, time.Time{}.UnixNano(), fl.LastModified.UnixNano())
	assertEqual(t, dt.UnixNano(), fl.LastUpdated.UnixNano())
	assertEqual(t, FetchResultOK, fl.Result)
	assertEqual(t, "message", fl.ResultMessage)
	assertEqual(t, dt.UnixNano(), fl.StartTime.UnixNano())
	assertEqual(t, 200, fl.StatusCode)
	assertEqual(t, "title", fl.Title)
	assertEqual(t, "url", fl.URL)
	assertEqual(t, false, fl.UsesGzip)
	assertEqual(t, 2, fl.NewEntries)

}
