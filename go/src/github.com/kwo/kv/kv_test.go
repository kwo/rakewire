package kv

import (
	"testing"
	"time"
)

func TestDates(t *testing.T) {

	t.Parallel()

	assertEqual := func(a interface{}, b interface{}) {
		if a != b {
			t.Errorf("Not equal: expected %v, actual %v", a, b)
		}
	}

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil || tz == nil {
		t.Fatalf("Timezone error: %v, %v\n", tz, err)
	}

	assertEqual("2015-11-20T20:42:55Z", dt.UTC().Format(time.RFC3339))
	assertEqual("2015-11-20T21:42:55+01:00", dt.In(tz).Format(time.RFC3339))

}

func TestEncode(t *testing.T) {

	t.Parallel()

	assertEqual := func(a interface{}, b interface{}) {
		if a != b {
			t.Errorf("Not equal: expected %v, actual %v", a, b)
		}
	}

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
	if err != nil || meta == nil || data == nil {
		t.Fatalf("Encode error: %v, %v, %v\n", meta, data, err)
	}

	assertEqual("object", meta.Name)
	assertEqual("ID", meta.Key)
	assertEqual(1, len(meta.Indexes))
	assertEqual(2, len(meta.Indexes["StartTime"]))
	assertEqual("StartTime", meta.Indexes["StartTime"][0])
	assertEqual("Key", meta.Indexes["StartTime"][1])

	assertEqual("object", data.Name)
	assertEqual("555", data.Key)
	assertEqual(1, len(data.Indexes))
	assertEqual(2, len(data.Indexes["StartTime"]))
	assertEqual("2015-11-20T20:42:55.000000033Z", data.Indexes["StartTime"][0])
	assertEqual("1", data.Indexes["StartTime"][1])

	assertEqual(4, len(data.Values))
	assertEqual("555", data.Values["ID"])
	assertEqual("1", data.Values["Key"])
	assertEqual("2015-11-20T20:42:55.000000033Z", data.Values["StartTime"])
	assertEqual("2015-11-20T20:42:55.000000033Z", data.Values["StartTime2"])

}

func TestEncodeNoKey(t *testing.T) {

	t.Parallel()

	type object struct {
		Key        int       `kv:"StartTime:2"`
		StartTime  time.Time `kv:"StartTime:1"`
		StartTime2 time.Time
	}

	o := &object{}

	if _, _, err := Encode(o); err == nil {
		t.Fatal("Encode error expected")
	} else if err.Error() != "Empty primary key for object" {
		t.Errorf("Wrong error message: %s", err.Error())
	}

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

	if _, _, err := Encode(o); err == nil {
		t.Fatal("Encode error expected")
	} else if err.Error() != "Non-contiguous index names for entity object, index StartTime" {
		t.Errorf("Wrong error message: %s", err.Error())
	}

}

func TestDecode(t *testing.T) {

	t.Parallel()

	assertEqual := func(a interface{}, b interface{}) {
		if a != b {
			t.Errorf("Not equal: expected %v, actual %v", a, b)
		}
	}

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
	if err != nil {
		t.Fatalf("Encode error: %v\n", err)
	}

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil || tz == nil {
		t.Fatalf("Timezone error: %v, %v\n", tz, err)
	}

	assertEqual("hello", o.ID)
	assertEqual(0, o.Key)
	assertEqual(dt, o.StartTime)
	assertEqual(dt, o.StartTime2.UTC())

}

func TestDecodeNonPointer(t *testing.T) {

	t.Parallel()

	type object struct {
	}

	o := object{}
	values := map[string]string{}

	if err := Decode(o, values); err == nil {
		t.Fatal("Decode error expected")
	} else if err.Error() != "Cannot decode non-pointer object" {
		t.Errorf("Wrong error message, expected '%s'", "Cannot decode non-pointer object")
	}

}

func TestDataFrom(t *testing.T) {

	t.Parallel()

	assertEqual := func(a interface{}, b interface{}) {
		if a != b {
			t.Errorf("Not equal: expected %v, actual %v", a, b)
		}
	}

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
	if data == nil {
		t.Fatal("Nil result from DataFrom")
	}

	assertEqual("hello", data.Key)
	assertEqual(1, len(data.Indexes))
	assertEqual("2015-11-20T20:42:55Z", data.Indexes["StartTime"][0])
	assertEqual("", data.Indexes["StartTime"][1])

}

func TestEncodeFields(t *testing.T) {

	t.Parallel()

	assertEqual := func(a interface{}, b interface{}) {
		if a != b {
			t.Errorf("Not equal: expected %v, actual %v", a, b)
		}
	}

	value, err := EncodeFields(1, 4.5, time.Date(2015, time.November, 20, 20, 42, 55, 0, time.UTC), "hello")
	if err != nil || value == nil {
		t.Fatalf("Nil result from EncodeFields: error: %v", err)
	}

	assertEqual(4, len(value))
	assertEqual("1", value[0])
	assertEqual("4.5", value[1])
	assertEqual("2015-11-20T20:42:55Z", value[2])
	assertEqual("hello", value[3])

}
