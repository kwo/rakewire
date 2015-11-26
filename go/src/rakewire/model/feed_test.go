package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/serial"
	"testing"
	"time"
)

func TestNewFeed(t *testing.T) {

	f := NewFeed("http://localhost/")
	require.NotNil(t, f)
	assert.Equal(t, "http://localhost/", f.URL)
	assert.NotNil(t, f.NextFetch)
	assert.NotNil(t, f.ID)
	assert.Equal(t, 0, len(f.Log))
	assert.Equal(t, 36, len(f.ID))

}

func TestFeedSerial(t *testing.T) {

	f := getNewFeed()
	validateFeed(t, f)

	meta, data, err := serial.Encode(f)
	require.Nil(t, err)
	require.NotNil(t, meta)
	require.NotNil(t, data)

	f2 := &Feed{}
	err = serial.Decode(f2, data.Values)
	require.Nil(t, err)
	validateFeed(t, f2)

}

func TestFeedJSON(t *testing.T) {

	f := getNewFeed()
	validateFeed(t, f)

	data, err := json.Marshal(f)
	require.Nil(t, err)
	require.NotNil(t, data)

	f2 := &Feed{}
	err = json.Unmarshal(data, f2)
	require.Nil(t, err)
	validateFeed(t, f2)

}

func TestFeedFunctions(t *testing.T) {

	f := NewFeed("http://localhost")

	assert.Nil(t, f.GetLast())

	f.AddLog(NewFeedLog(f.ID))
	f.AddLog(NewFeedLog(f.ID))
	f.AddLog(NewFeedLog(f.ID))
	f.Last200 = f.AddLog(NewFeedLog(f.ID))
	f.Last = f.AddLog(NewFeedLog(f.ID))

	assert.Equal(t, 5, len(f.Log))
	assert.Equal(t, f.Last, f.GetLast().ID)
	assert.Equal(t, f.Last200, f.GetLast200().ID)

}

func TestAdjustFetchTime(t *testing.T) {

	f := NewFeed("http://localhost")
	require.NotNil(t, f)
	require.NotNil(t, f.NextFetch)
	assert.False(t, f.NextFetch.IsZero()) // nextfetch is initialized to now

	now := time.Now()
	f.NextFetch = now

	diff := 3 * time.Hour
	f.AdjustFetchTime(diff)
	assert.Equal(t, now.Add(diff).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

}

func TestUpdateFetchTime(t *testing.T) {

	f := NewFeed("http://localhost")
	require.NotNil(t, f)
	require.NotNil(t, f.NextFetch)
	assert.False(t, f.NextFetch.IsZero()) // nextfetch is initialized to now

	now := time.Now()

	f.UpdateFetchTime(now.Add(-29 * time.Minute))
	assert.Equal(t, now.Add(10*time.Minute).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-30 * time.Minute))
	assert.Equal(t, now.Add(1*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-3 * time.Hour))
	assert.Equal(t, now.Add(1*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-72 * time.Hour))
	assert.Equal(t, now.Add(24*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(time.Time{})
	assert.Equal(t, now.Add(24*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

}

func getNewFeed() *Feed {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	f := NewFeed("http://localhost")
	f.Last = "last"
	f.Last200 = "last200"
	f.LastUpdated = dt
	f.NextFetch = dt
	f.Notes = "notes"
	f.Title = "title"

	return f

}

func validateFeed(t *testing.T, f *Feed) {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	require.NotNil(t, f)
	assert.Equal(t, "last", f.Last)
	assert.Equal(t, "last200", f.Last200)
	assert.Equal(t, dt.UnixNano(), f.LastUpdated.UnixNano())
	assert.Equal(t, dt.UnixNano(), f.NextFetch.UnixNano())
	assert.Equal(t, "notes", f.Notes)
	assert.Equal(t, "title", f.Title)

}
