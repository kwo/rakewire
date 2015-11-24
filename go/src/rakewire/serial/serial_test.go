package serial

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/logging"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	// initialize logging
	logging.Init(&logging.Configuration{
		Level: "debug",
	})
	logger.Debug("Logging configured")

	m.Run()

}

func TestDates(t *testing.T) {

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	require.Nil(t, err)
	require.NotNil(t, tz)

	assert.Equal(t, "2015-11-20T20:42:55Z", dt.UTC().Format(time.RFC3339))
	assert.Equal(t, "2015-11-20T21:42:55+01:00", dt.In(tz).Format(time.RFC3339))

}

func TestEncode(t *testing.T) {

	type object struct {
		ID         string
		Key        int       `db:"StartTime:2"`
		StartTime  time.Time `db:"StartTime:1"`
		StartTime2 time.Time
		Log        int64 `db:"-"`
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
	require.Nil(t, err)
	require.NotNil(t, meta)
	require.NotNil(t, data)

	assert.Equal(t, "object", meta.Name)
	assert.Equal(t, "ID", meta.Key)
	assert.Equal(t, 1, len(meta.Indexes))
	assert.Equal(t, 2, len(meta.Indexes["StartTime"]))
	assert.Equal(t, "StartTime", meta.Indexes["StartTime"][0])
	assert.Equal(t, "Key", meta.Indexes["StartTime"][1])

	assert.Equal(t, "object", data.Name)
	assert.Equal(t, "555", data.Key)
	assert.Equal(t, 1, len(data.Indexes))
	assert.Equal(t, 2, len(data.Indexes["StartTime"]))
	assert.Equal(t, "2015-11-20T20:42:55.000000033Z", data.Indexes["StartTime"][0])
	assert.Equal(t, "1", data.Indexes["StartTime"][1])

	assert.Equal(t, 4, len(data.Values))
	assert.Equal(t, "555", data.Values["ID"])
	assert.Equal(t, "1", data.Values["Key"])
	assert.Equal(t, "2015-11-20T20:42:55.000000033Z", data.Values["StartTime"])
	assert.Equal(t, "2015-11-20T20:42:55.000000033Z", data.Values["StartTime2"])

}

func TestEncodeNoKey(t *testing.T) {

	type object struct {
		Key        int       `db:"StartTime:2"`
		StartTime  time.Time `db:"StartTime:1"`
		StartTime2 time.Time
	}

	o := &object{}

	_, _, err := Encode(o)
	require.NotNil(t, err)
	assert.Equal(t, "Empty primary key for object.", err.Error())

}

func TestEncodeNonContiguousIndexes(t *testing.T) {

	type object struct {
		ID         string
		Key        int       `db:"StartTime:3"`
		StartTime  time.Time `db:"StartTime:1"`
		StartTime2 time.Time
	}

	o := &object{}

	_, _, err := Encode(o)
	require.NotNil(t, err)
	assert.Equal(t, "Non-contiguous index names for entity object, index StartTime.", err.Error())

}

func TestDecode(t *testing.T) {

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
	require.Nil(t, err)

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 33, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	require.Nil(t, err)
	require.NotNil(t, tz)

	assert.Equal(t, "hello", o.ID)
	assert.Equal(t, 0, o.Key)
	assert.Equal(t, dt, o.StartTime)
	assert.Equal(t, dt, o.StartTime2.UTC())

}

func TestDecodeNonPointer(t *testing.T) {

	type object struct {
	}

	o := object{}
	values := map[string]string{}

	err := Decode(o, values)
	require.NotNil(t, err)
	assert.Equal(t, "Cannot decode non-pointer object", err.Error())

}

func TestDataFrom(t *testing.T) {

	type object struct {
		ID        string
		Key       int       `db:"StartTime:2"`
		StartTime time.Time `db:"StartTime:1"`
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
	require.NotNil(t, data)

	assert.Equal(t, "hello", data.Key)
	assert.Equal(t, 1, len(data.Indexes))
	assert.Equal(t, "2015-11-20T20:42:55Z", data.Indexes["StartTime"][0])
	assert.Equal(t, "", data.Indexes["StartTime"][1])

}

func TestEncodeFields(t *testing.T) {

	value, err := EncodeFields(1, 4.5, time.Date(2015, time.November, 20, 20, 42, 55, 0, time.UTC), "hello")
	require.Nil(t, err)
	require.NotNil(t, value)
	assert.Equal(t, 4, len(value))

	assert.Equal(t, "1", value[0])
	assert.Equal(t, "4.5", value[1])
	assert.Equal(t, "2015-11-20T20:42:55Z", value[2])
	assert.Equal(t, "hello", value[3])

}
