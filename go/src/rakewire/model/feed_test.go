package model

import (
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
	if f.ID != empty {
		t.Errorf("f.ID not equal, expected: %d, actual: %d", 0, f.ID)
	}

}

func TestFeedSerial(t *testing.T) {

	t.Parallel()

	f := getNewFeed()
	f.ID = kvKeyUintEncode(1)
	validateFeed(t, f)

	data := f.serialize()
	assertNotNil(t, data)

	f2 := &Feed{}
	err := f2.deserialize(data)
	assertNoError(t, err)
	validateFeed(t, f2)

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

	expectedTime := f.LastUpdated.Add(15 * time.Minute)
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

	expectedTime := f.LastUpdated.Add(45 * time.Minute)
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

	expectedTime := f.LastUpdated.Add(4 * time.Hour)
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

	expectedTime := f.LastUpdated.Add(((4 * 24) + 1) * time.Hour)
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
