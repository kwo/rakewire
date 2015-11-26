package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/serial"
	"testing"
	"time"
)

func TestNewFeedLog(t *testing.T) {

	fl := NewFeedLog("123")
	require.NotNil(t, fl)
	assert.Equal(t, "123", fl.FeedID)
	assert.NotNil(t, fl.ID)
	assert.Equal(t, 36, len(fl.ID))

}

func TestFeedLogSerial(t *testing.T) {

	fl := getNewFeedLog()
	validateFeedLog(t, fl)

	meta, data, err := serial.Encode(fl)
	require.Nil(t, err)
	require.NotNil(t, meta)
	require.NotNil(t, data)

	fl2 := &FeedLog{}
	err = serial.Decode(fl2, data.Values)
	require.Nil(t, err)
	validateFeedLog(t, fl2)

}

func TestFeedLogJSON(t *testing.T) {

	fl := getNewFeedLog()
	validateFeedLog(t, fl)

	data, err := json.Marshal(fl)
	require.Nil(t, err)
	require.NotNil(t, data)

	fl2 := &FeedLog{}
	err = json.Unmarshal(data, fl2)
	require.Nil(t, err)
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

	require.NotNil(t, fl)
	assert.Equal(t, 0, fl.ContentLength)
	assert.Equal(t, "text/plain", fl.ContentType)
	assert.Equal(t, 6*time.Second, fl.Duration)
	assert.Equal(t, "etag", fl.ETag)
	assert.Equal(t, "123", fl.FeedID)
	assert.Equal(t, "flavor", fl.Flavor)
	assert.Equal(t, "", fl.Generator)
	assert.Equal(t, 36, len(fl.ID))
	assert.Equal(t, true, fl.IsUpdated)
	assert.Equal(t, time.Time{}.UnixNano(), fl.LastModified.UnixNano())
	assert.Equal(t, dt.UnixNano(), fl.LastUpdated.UnixNano())
	assert.Equal(t, FetchResultOK, fl.Result)
	assert.Equal(t, "message", fl.ResultMessage)
	assert.Equal(t, dt.UnixNano(), fl.StartTime.UnixNano())
	assert.Equal(t, 200, fl.StatusCode)
	assert.Equal(t, "title", fl.Title)
	assert.Equal(t, "url", fl.URL)
	assert.Equal(t, UpdateCheckFeed, fl.UpdateCheck)
	assert.Equal(t, false, fl.UsesGzip)

}
