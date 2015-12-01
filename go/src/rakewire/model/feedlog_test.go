package model

import (
	"encoding/json"
	"rakewire/kv"
	"testing"
	"time"
)

func TestNewFeedLog(t *testing.T) {

	t.Parallel()

	fl := NewFeedLog("123")
	assertNotNil(t, fl)
	assertEqual(t, "123", fl.FeedID)
	assertNotNil(t, fl.ID)
	assertEqual(t, 36, len(fl.ID))

}

func TestFeedLogkv(t *testing.T) {

	t.Parallel()

	fl := getNewFeedLog()
	validateFeedLog(t, fl)

	meta, data, err := kv.Encode(fl)
	assertNoError(t, err)
	assertNotNil(t, meta)
	assertNotNil(t, data)

	fl2 := &FeedLog{}
	err = kv.Decode(fl2, data.Values)
	assertNoError(t, err)
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
	fl.IsUpdated = true
	//fl.LastModified = dt
	fl.LastUpdated = dt
	fl.Result = FetchResultOK
	fl.ResultMessage = "message"
	fl.StartTime = dt
	fl.StatusCode = 200
	fl.Title = "title"
	fl.URL = "url"
	fl.UpdateCheck = UpdateCheckFeed
	fl.UsesGzip = false

	return fl

}

func validateFeedLog(t *testing.T, fl *FeedLog) {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	assertNotNil(t, fl)
	assertEqual(t, 0, fl.ContentLength)
	assertEqual(t, "text/plain", fl.ContentType)
	assertEqual(t, 6*time.Second, fl.Duration)
	assertEqual(t, "etag", fl.ETag)
	assertEqual(t, "123", fl.FeedID)
	assertEqual(t, "flavor", fl.Flavor)
	assertEqual(t, "", fl.Generator)
	assertEqual(t, 36, len(fl.ID))
	assertEqual(t, true, fl.IsUpdated)
	assertEqual(t, time.Time{}.UnixNano(), fl.LastModified.UnixNano())
	assertEqual(t, dt.UnixNano(), fl.LastUpdated.UnixNano())
	assertEqual(t, FetchResultOK, fl.Result)
	assertEqual(t, "message", fl.ResultMessage)
	assertEqual(t, dt.UnixNano(), fl.StartTime.UnixNano())
	assertEqual(t, 200, fl.StatusCode)
	assertEqual(t, "title", fl.Title)
	assertEqual(t, "url", fl.URL)
	assertEqual(t, UpdateCheckFeed, fl.UpdateCheck)
	assertEqual(t, false, fl.UsesGzip)

}
