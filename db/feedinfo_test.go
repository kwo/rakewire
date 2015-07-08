package db

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {

	dt := time.Date(2015, time.July, 8, 9, 24, 0, 0, time.UTC)

	fi, err := NewFeedInfo("http://localhost:8888/")
	require.Nil(t, err)
	require.NotNil(t, fi)
	fi.LastUpdated = &dt

	assert.Nil(t, fi.LastFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.Equal(t, "http://localhost:8888/", fi.URL)
	assert.EqualValues(t, dt, *fi.LastUpdated)

	data, err := fi.Marshal()
	require.Nil(t, err)
	require.NotNil(t, data)
	assert.Equal(t, 124, len(data))

	//fmt.Println(string(data))
	//fmt.Printf("size: %d\n", len(data))

	fi2 := FeedInfo{}
	err = fi2.Unmarshal(data)
	require.Nil(t, err)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.Nil(t, fi2.LastFetch)
	assert.NotNil(t, fi2.LastUpdated)
	assert.EqualValues(t, dt, *fi2.LastUpdated)

}
