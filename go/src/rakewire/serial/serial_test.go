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

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 0, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	require.Nil(t, err)
	require.NotNil(t, tz)

	assert.Equal(t, "2015-11-20T20:42:55Z", dt.UTC().Format(time.RFC3339))
	assert.Equal(t, "2015-11-20T21:42:55+01:00", dt.In(tz).Format(time.RFC3339))

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
		"StartTime":  "2015-11-20T20:42:55Z",
		"StartTime2": "2015-11-20T21:42:55+01:00",
	}

	err := Decode(o, values)
	require.Nil(t, err)

	dt := time.Date(2015, time.November, 20, 20, 42, 55, 0, time.UTC)
	tz, err := time.LoadLocation("Europe/Berlin")
	require.Nil(t, err)
	require.NotNil(t, tz)

	assert.Equal(t, "hello", o.ID)
	assert.Equal(t, 0, o.Key)
	assert.Equal(t, dt, o.StartTime)
	assert.Equal(t, dt, o.StartTime2.UTC())

}

func TestDataFrom(t *testing.T) {

	type object struct {
		ID        string
		Key       int       `db:"indexStartTime:2"`
		StartTime time.Time `db:"indexStartTime:1"`
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
