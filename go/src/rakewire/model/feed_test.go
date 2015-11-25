package model

import (
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

func TestFeed(t *testing.T) {

	dt := time.Date(2015, time.July, 8, 9, 24, 0, 0, time.UTC)

	fi := NewFeed("http://localhost:8888/")
	require.NotNil(t, fi)
	fi.LastUpdated = dt

	assert.NotNil(t, fi.NextFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.Equal(t, "http://localhost:8888/", fi.URL)
	assert.EqualValues(t, dt, fi.LastUpdated)

	meta, data, err := serial.Encode(fi)
	require.Nil(t, err)
	require.NotNil(t, meta)
	require.NotNil(t, data)
	//assert.Equal(t, 190, len(data))

	// fmt.Println(string(data))
	// fmt.Printf("size: %d\n", len(data))

	fi2 := &Feed{}
	err = serial.Decode(fi2, data.Values)
	require.Nil(t, err)
	assert.NotNil(t, fi2)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.NotNil(t, fi2.NextFetch)
	assert.NotNil(t, fi2.LastUpdated)
	assert.EqualValues(t, dt, fi2.LastUpdated)

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
