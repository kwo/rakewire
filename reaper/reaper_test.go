package reaper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	"rakewire.com/db/bolt"
	"testing"
)

const (
	databaseFile = "../test/pollfeed.db"
)

func TestReaper(t *testing.T) {

	t.SkipNow()

	// open database
	database := &bolt.Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// create service
	cfg := &Configuration{}
	pf := NewService(cfg, database)

	pf.Start()
	require.Equal(t, true, pf.IsRunning())
	pf.Stop()
	assert.Equal(t, false, pf.IsRunning())

	// close database
	err = database.Close()
	assert.Nil(t, err)

	// remove file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
