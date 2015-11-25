package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"rakewire/db"
	"rakewire/logging"
	"testing"
)

const (
	feedFile              = "../../../../test/feedlist.txt"
	databaseTempDirectory = "../../../../test"
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

func getTempFile(t *testing.T) string {
	f, err := ioutil.TempFile(empty, "bolt-")
	require.Nil(t, err)
	f.Close()
	return f.Name()
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
