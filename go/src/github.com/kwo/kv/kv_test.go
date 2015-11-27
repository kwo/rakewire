package kv

import (
	"testing"
	"time"
)

func TestDates(t *testing.T) {

	t.Parallel()

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	assertNoError(t, err)
	assertNotNil(t, tz)

	assertEqual(t, "2015-11-20T20:42:55Z", dt.UTC().Format(time.RFC3339))
	assertEqual(t, "2015-11-20T21:42:55+01:00", dt.In(tz).Format(time.RFC3339))

}

func TestEncode(t *testing.T) {

	t.Parallel()

	type object struct {
		ID         string
		Key        int       `kv:"StartTime:2"`
		StartTime  time.Time `kv:"StartTime:1"`
		StartTime2 time.Time
		Log        int64 `kv:"-"`
	}

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)

	o := &object{
		ID:         "555",
		Key:        1,
		StartTime:  dt,
		StartTime2: dt.Local(),
		Log:        44,
	}

	meta, data, err := Encode(o)
	assertNoError(t, err)
	assertNotNil(t, meta)
	assertNotNil(t, data)

	assertEqual(t, "object", meta.Name)
	assertEqual(t, "ID", meta.Key)
	assertEqual(t, 1, len(meta.Indexes))
	assertEqual(t, 2, len(meta.Indexes["StartTime"]))
	assertEqual(t, "StartTime", meta.Indexes["StartTime"][0])
	assertEqual(t, "Key", meta.Indexes["StartTime"][1])

	assertEqual(t, "object", data.Name)
	assertEqual(t, "555", data.Key)
	assertEqual(t, 1, len(data.Indexes))
	assertEqual(t, 2, len(data.Indexes["StartTime"]))
	assertEqual(t, "2015-11-20T20:42:55.000000033Z", data.Indexes["StartTime"][0])
	assertEqual(t, "1", data.Indexes["StartTime"][1])

	assertEqual(t, 4, len(data.Values))
	assertEqual(t, "555", data.Values["ID"])
	assertEqual(t, "1", data.Values["Key"])
	assertEqual(t, "2015-11-20T20:42:55.000000033Z", data.Values["StartTime"])
	assertEqual(t, "2015-11-20T20:42:55.000000033Z", data.Values["StartTime2"])

}

func TestEncodeNoKey(t *testing.T) {

	t.Parallel()

	type object struct {
		Key        int       `kv:"StartTime:2"`
		StartTime  time.Time `kv:"StartTime:1"`
		StartTime2 time.Time
	}

	o := &object{}

	_, _, err := Encode(o)
	assertNotNil(t, err)
	assertEqual(t, "Empty primary key for object.", err.Error())

}

func TestEncodeNonContiguousIndexes(t *testing.T) {

	t.Parallel()

	type object struct {
		ID         string
		Key        int       `kv:"StartTime:3"`
		StartTime  time.Time `kv:"StartTime:1"`
		StartTime2 time.Time
	}

	o := &object{}

	_, _, err := Encode(o)
	assertNotNil(t, err)
	assertEqual(t, "Non-contiguous index names for entity object, index StartTime.", err.Error())

}

func TestDecode(t *testing.T) {

	t.Parallel()

	type object struct {
		ID         string
		Key        int
		StartTime  time.Time
		StartTime2 time.Time
	}

	o := &object{}
	values := map[string]string{
		"ID":         "hello",
		"StartTime":  "2015-11-20T20:42:55.000000033Z",
		"StartTime2": "2015-11-20T21:42:55.000000033+01:00",
	}

	err := Decode(o, values)
	assertNoError(t, err)

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	assertNoError(t, err)
	assertNotNil(t, tz)

	assertEqual(t, "hello", o.ID)
	assertEqual(t, 0, o.Key)
	assertEqual(t, dt, o.StartTime)
	assertEqual(t, dt, o.StartTime2.UTC())

}

func TestDecodeNonPointer(t *testing.T) {

	t.Parallel()

	type object struct {
	}

	o := object{}
	values := map[string]string{}

	err := Decode(o, values)
	assertNotNil(t, err)
	assertEqual(t, "Cannot decode non-pointer object", err.Error())

}

func TestDataFrom(t *testing.T) {

	t.Parallel()

	type object struct {
		ID        string
		Key       int       `kv:"StartTime:2"`
		StartTime time.Time `kv:"StartTime:1"`
	}

	values := map[string]string{
		"ID":        "hello",
		"StartTime": "2015-11-20T20:42:55Z",
	}

	metadata := &Metadata{
		Name: "object",
		Key:  "ID",
		Indexes: map[string][]string{
			"StartTime": {"StartTime", "Key"},
		},
	}

	data := DataFrom(metadata, values)
	assertNotNil(t, data)

	assertEqual(t, "hello", data.Key)
	assertEqual(t, 1, len(data.Indexes))
	assertEqual(t, "2015-11-20T20:42:55Z", data.Indexes["StartTime"][0])
	assertEqual(t, "", data.Indexes["StartTime"][1])

}

func TestEncodeFields(t *testing.T) {

	t.Parallel()

	value, err := EncodeFields(1, 4.5, time.Date(2015, time.November, 20, 20, 42, 55, 0, time.UTC), "hello")
	assertNoError(t, err)
	assertNotNil(t, value)

	assertEqual(t, 4, len(value))
	assertEqual(t, "1", value[0])
	assertEqual(t, "4.5", value[1])
	assertEqual(t, "2015-11-20T20:42:55Z", value[2])
	assertEqual(t, "hello", value[3])

}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Fatal("Expected nil value")
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("Expected not nil value")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Not equal: expected %v, actual %v", a, b)
	}
}
