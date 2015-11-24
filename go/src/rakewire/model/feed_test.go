package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/serial"
	"testing"
	"time"
)

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

func TestFeeds(t *testing.T) {

	dt := time.Date(2015, time.July, 8, 9, 24, 0, 0, time.UTC)

	fi := NewFeed("http://localhost:8888/")
	require.NotNil(t, fi)
	fi.LastUpdated = dt

	assert.NotNil(t, fi.NextFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.Equal(t, "http://localhost:8888/", fi.URL)
	assert.EqualValues(t, dt, fi.LastUpdated)

}
