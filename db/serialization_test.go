package db

import (
	"encoding/json"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJson(t *testing.T) {

	fi := FeedInfo{
		ID:  "ABCDEFG",
		URL: "http://localhost:8888/",
	}

	data, err := json.Marshal(fi)
	require.Nil(t, err)
	require.NotNil(t, data)
	assert.True(t, fi.LastFetch.IsZero())

	//fmt.Println(string(data))

	var fi2 FeedInfo
	err = json.Unmarshal(data, &fi2)
	require.Nil(t, err)
	assert.Equal(t, fi.ID, fi2.ID)
	assert.Equal(t, fi.URL, fi2.URL)
	assert.True(t, fi2.LastFetch.IsZero())

	// times are deserialized with a non-nil location
	// assert.EqualValues(t, fi, fi2)

}
