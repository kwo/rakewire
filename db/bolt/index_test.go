package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire.com/config"
	"testing"
	"time"
)

func TestIndexFetchNext(t *testing.T) {

	//t.SkipNow()

	cfg := config.GetConfig()

	// open database
	database := &Database{}
	err := database.Open(&cfg.Database)
	require.Nil(t, err)

	feeds, err := database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	maxTime := time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	// close database
	err = database.Close()
	assert.Nil(t, err)

}
