package bolt

import (
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
	require.NotNil(t, d)

}

func openDatabase(t *testing.T) *Database {

	f, err := ioutil.TempFile(empty, "bolt-")
	require.Nil(t, err)
	filename := f.Name()
	f.Close()

	database := &Database{}
	err = database.Open(&db.Configuration{
		Location: filename,
	})
	require.Nil(t, err)

	return database

}

func closeDatabase(t *testing.T, database *Database) {

	// close database
	err := database.Close()
	require.Nil(t, err)
	require.Nil(t, database.db)

	// remove file
	err = os.Remove(database.databaseFile)
	require.Nil(t, err)

}
