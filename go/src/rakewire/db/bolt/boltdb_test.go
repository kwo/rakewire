package bolt

import (
	"io/ioutil"
	"os"
	"rakewire/db"
	"testing"
)

const (
	feedFile              = "../../../../test/feedlist.txt"
	databaseTempDirectory = "../../../../test"
)

func TestInterface(t *testing.T) {

	var d db.Database = &Database{}
	assertNotNil(t, d)

}

func openDatabase(t *testing.T) *Database {

	f, err := ioutil.TempFile(empty, "bolt-")
	assertNoError(t, err)
	filename := f.Name()
	f.Close()

	database := NewService(&db.Configuration{
		Location: filename,
	})
	err = database.Open()
	assertNoError(t, err)
	if !database.running {
		t.Error("database is not running")
	}

	return database

}

func closeDatabase(t *testing.T, database *Database) {

	// close database
	database.Close()
	if database.running {
		t.Error("database is still running")
	}
	if database.db != nil {
		t.Error("database.db is not nil")
	}

	// remove file
	err := os.Remove(database.databaseFile)
	assertNoError(t, err)

}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("Expected not nil value")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Not equal: expected %v, actual %v", a, b)
	}
}
