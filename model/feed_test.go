package model

import (
	"bytes"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFeed(t *testing.T) {

	dt := time.Date(2015, time.July, 8, 9, 24, 0, 0, time.UTC)

	fi := NewFeed("http://localhost:8888/")
	require.NotNil(t, fi)
	fi.LastUpdated = &dt

	assert.Nil(t, fi.LastFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.Equal(t, "http://localhost:8888/", fi.URL)
	assert.EqualValues(t, dt, *fi.LastUpdated)

	data, err := fi.Encode()
	require.Nil(t, err)
	require.NotNil(t, data)
	assert.Equal(t, 137, len(data))

	// fmt.Println(string(data))
	// fmt.Printf("size: %d\n", len(data))

	fi2 := Feed{}
	err = fi2.Decode(data)
	require.Nil(t, err)
	assert.NotNil(t, fi2)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.Nil(t, fi2.LastFetch)
	assert.NotNil(t, fi2.LastUpdated)
	assert.EqualValues(t, dt, *fi2.LastUpdated)

}

func TestFeeds(t *testing.T) {

	dt := time.Date(2015, time.July, 8, 9, 24, 0, 0, time.UTC)

	fi := NewFeed("http://localhost:8888/")
	require.NotNil(t, fi)
	fi.LastUpdated = &dt

	assert.Nil(t, fi.LastFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.Equal(t, "http://localhost:8888/", fi.URL)
	assert.EqualValues(t, dt, *fi.LastUpdated)

	fds := NewFeeds()
	var buf bytes.Buffer

	fds.Add(fi)
	assert.Equal(t, 1, fds.Size())
	err := fds.Serialize(&buf)
	data := buf.Bytes()
	require.Nil(t, err)
	require.NotNil(t, data)
	assert.Equal(t, 127, len(data))

	// fmt.Println(string(data))
	// fmt.Printf("size: %d\n", len(data))

	fds = NewFeeds()
	assert.Equal(t, 0, fds.Size())
	err = fds.Deserialize(bytes.NewReader(data))
	require.Nil(t, err)
	assert.Equal(t, 1, fds.Size())
	fi2 := fds.Get(fi.ID)
	assert.NotNil(t, fi2)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.Nil(t, fi2.LastFetch)
	assert.NotNil(t, fi2.LastUpdated)
	assert.EqualValues(t, dt, *fi2.LastUpdated)

}
