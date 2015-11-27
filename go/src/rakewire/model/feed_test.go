package model

import (
	"encoding/json"
	"github.com/kwo/kv"
	"testing"
	"time"
)

func TestNewFeed(t *testing.T) {

	t.Parallel()

	f := NewFeed("http://localhost/")
	assertNotNil(t, f)
	assertEqual(t, "http://localhost/", f.URL)
	assertNotNil(t, f.NextFetch)
	assertNotNil(t, f.ID)
	assertEqual(t, 0, len(f.Log))
	assertEqual(t, 36, len(f.ID))

}

func TestFeedSerial(t *testing.T) {

	t.Parallel()

	f := getNewFeed()
	validateFeed(t, f)

	meta, data, err := kv.Encode(f)
	assertNoError(t, err)
	assertNotNil(t, meta)
	assertNotNil(t, data)

	f2 := &Feed{}
	err = kv.Decode(f2, data.Values)
	assertNoError(t, err)
	validateFeed(t, f2)

}

func TestFeedJSON(t *testing.T) {

	t.Parallel()

	f := getNewFeed()
	validateFeed(t, f)

	data, err := json.Marshal(f)
	assertNoError(t, err)
	assertNotNil(t, data)

	f2 := &Feed{}
	err = json.Unmarshal(data, f2)
	assertNoError(t, err)
	validateFeed(t, f2)

}

func TestFeedsJSON(t *testing.T) {

	t.Parallel()

	f := getNewFeed()
	validateFeed(t, f)

	var feeds []*Feed
	feeds = append(feeds, f)

	data, err := json.Marshal(&feeds)
	assertNoError(t, err)
	assertNotNil(t, data)

	var feeds2 []*Feed
	err = json.Unmarshal(data, &feeds2)
	assertNoError(t, err)
	assertEqual(t, 1, len(feeds2))
	validateFeed(t, feeds2[0])

}

func TestFeedFunctions(t *testing.T) {

	t.Parallel()

	f := NewFeed("http://localhost")

	f.AddLog(NewFeedLog(f.ID))
	f.AddLog(NewFeedLog(f.ID))
	f.AddLog(NewFeedLog(f.ID))

	assertEqual(t, 3, len(f.Log))

}

func TestAdjustFetchTime(t *testing.T) {

	t.Parallel()

	f := NewFeed("http://localhost")
	assertNotNil(t, f)
	assertNotNil(t, f.NextFetch)
	assertEqual(t, false, f.NextFetch.IsZero()) // nextfetch is initialized to now

	now := time.Now()
	f.NextFetch = now

	diff := 3 * time.Hour
	f.AdjustFetchTime(diff)
	assertEqual(t, now.Add(diff).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

}

func TestUpdateFetchTime(t *testing.T) {

	t.Parallel()

	f := NewFeed("http://localhost")
	assertNotNil(t, f)
	assertNotNil(t, f.NextFetch)
	assertEqual(t, false, f.NextFetch.IsZero()) // nextfetch is initialized to now

	now := time.Now()

	f.UpdateFetchTime(now.Add(-29 * time.Minute))
	assertEqual(t, now.Add(10*time.Minute).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-30 * time.Minute))
	assertEqual(t, now.Add(1*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-3 * time.Hour))
	assertEqual(t, now.Add(1*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(now.Add(-72 * time.Hour))
	assertEqual(t, now.Add(24*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

	f.UpdateFetchTime(time.Time{})
	assertEqual(t, now.Add(24*time.Hour).Truncate(time.Millisecond), f.NextFetch.Truncate(time.Millisecond))

}

func getNewFeed() *Feed {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	f := NewFeed("http://localhost")
	f.LastUpdated = dt
	f.NextFetch = dt
	f.Notes = "notes"
	f.Title = "title"

	return f

}

func validateFeed(t *testing.T, f *Feed) {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	assertNotNil(t, f)
	assertEqual(t, dt.UnixNano(), f.LastUpdated.UnixNano())
	assertEqual(t, dt.UnixNano(), f.NextFetch.UnixNano())
	assertEqual(t, "notes", f.Notes)
	assertEqual(t, "title", f.Title)

}
