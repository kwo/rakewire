package bolt

import (
	"github.com/stretchr/testify/assert"
	"os"
	"rakewire/db"
	"rakewire/logging"
	"testing"
)

const (
	feedFile     = "../../../../test/feedlist.txt"
	databaseFile = "../../../../test/bolt.db"
)

func TestMain(m *testing.M) {

	// initialize logging
	logging.Init(&logging.Configuration{
		Level: "debug",
	})

	logger.Debug("Logging configured")

	m.Run()

}

func TestInterface(t *testing.T) {

	var d db.Database = &Database{}
	assert.NotNil(t, d)

}

func cleanUp(t *testing.T, database *Database) {

	// close database
	err := database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	// remove file
	err = os.Remove(database.databaseFile)
	assert.Nil(t, err)

}
