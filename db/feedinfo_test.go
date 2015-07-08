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

	fi := FeedInfo{
		ID:          "ABCDEFG",
		URL:         "http://localhost:8888/",
		LastUpdated: &dt,
	}

	data, err := fi.Marshal()
	require.Nil(t, err)
	require.NotNil(t, data)
	assert.Nil(t, fi.LastFetch)
	assert.NotNil(t, fi.LastUpdated)
	assert.EqualValues(t, dt, *fi.LastUpdated)
	assert.Equal(t, 111, len(data))

	// jsonData, _ := unzip(data)
	// fmt.Println(string(jsonData))
	// fmt.Printf("size: %d\n", len(data))

	fi2 := FeedInfo{}
	err = fi2.Unmarshal(data)
	require.Nil(t, err)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.Nil(t, fi2.LastFetch)
	assert.NotNil(t, fi2.LastUpdated)
	assert.EqualValues(t, dt, *fi2.LastUpdated)

	// times are deserialized with a non-nil location
	// assert.EqualValues(t, fi, fi2)

}
