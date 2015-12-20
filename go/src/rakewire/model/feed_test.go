package model

import (
	"encoding/json"
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
	if f.ID != 0 {
		t.Errorf("f.ID not equal, expected: %d, actual: %d", 0, f.ID)
	}

}

func TestFeedSerial(t *testing.T) {

	t.Parallel()

	f := getNewFeed()
	validateFeed(t, f)

	data := f.Serialize()
	assertNotNil(t, data)

	f2 := &Feed{}
	err := f2.Deserialize(data)
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
	assertEqual(t, now.Add(diff).Truncate(time.Second), f.NextFetch.Truncate(time.Second))

}

func TestUpdateFetchTime1(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := NewFeed("http://localhost")
	if f == nil {
		t.Fatal("NewFeed returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("NewFeed must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-47 * time.Second)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(10 * time.Minute)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime2(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := NewFeed("http://localhost")
	if f == nil {
		t.Fatal("NewFeed returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("NewFeed must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-29 * time.Minute)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(40 * time.Minute)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime3(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := NewFeed("http://localhost")
	if f == nil {
		t.Fatal("NewFeed returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("NewFeed must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-3 * time.Hour).Add(-47 * time.Minute)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(5 * time.Hour)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime4(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := NewFeed("http://localhost")
	if f == nil {
		t.Fatal("NewFeed returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("NewFeed must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-4 * 24 * time.Hour)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(5 * 24 * time.Hour)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

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
