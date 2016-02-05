package model

import (
	"testing"
	"time"
)

func TestNewTransmission(t *testing.T) {

	t.Parallel()

	fl := NewTransmission(123)

	if fl == nil {
		t.Fatalf("Transmission factory not returning a valid transmission")
	}

	if fl.ID != 0 {
		t.Errorf("Factory method should not set th ID, expected %d, actual %d", 0, fl.ID)
	}

	if fl.FeedID != 123 {
		t.Errorf("Factory method not setting FeedID properly, expected %s, actual %s", "123", fl.FeedID)
	}

}

func TestTransmissionSerialize(t *testing.T) {

	t.Parallel()

	fl := getNewTransmission()
	validateTransmission(t, fl)

	data := fl.serialize()
	if data == nil {
		t.Fatal("Transmission serialize returned a nil map")
	}

	fl2 := &Transmission{}
	if err := fl2.deserialize(data); err != nil {
		t.Fatalf("Transmission deserialize returned an error: %s", err.Error())
	}
	validateTransmission(t, fl2)

}

func getNewTransmission() *Transmission {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	fl := NewTransmission(123)
	fl.ID = 1
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
	fl.NewItems = 2

	return fl

}

func validateTransmission(t *testing.T, fl *Transmission) {

	dt := time.Date(2015, time.November, 26, 13, 55, 0, 0, time.Local)

	assertNotNil(t, fl)

	if fl.ID != 1 {
		t.Errorf("Transmission ID incorrect, expected %d, actual %d", 0, fl.ID)
	}

	assertEqual(t, 0, fl.ContentLength)
	assertEqual(t, "text/plain", fl.ContentType)
	assertEqual(t, 6*time.Second, fl.Duration)
	assertEqual(t, "etag", fl.ETag)
	if fl.FeedID != 123 {
		t.Errorf("FeedIDs do not match, expected %d, actual %d", 123, fl.FeedID)
	}
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
	assertEqual(t, 2, fl.NewItems)

}
