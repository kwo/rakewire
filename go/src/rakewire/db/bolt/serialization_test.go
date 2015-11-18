package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire/db"
	m "rakewire/model"
	"testing"
)

func TestSerialization(t *testing.T) {

	// init db

	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// start testing

	database.Lock()
	err = database.db.Update(func(tx *bolt.Tx) error {
		fl := m.NewFeedLog("myFeedID")
		return marshal(fl, tx)
	})
	database.Unlock()
	assert.Nil(t, err)

	// cleanup

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
