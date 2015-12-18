package bolt

import (
	"io/ioutil"
	"os"
	"rakewire/db"
	m "rakewire/model"
	"testing"
)

const (
	feedFile              = "../../../../test/feedlist.txt"
	databaseTempDirectory = "../../../../test"
)

func TestInterfaceDatabase(t *testing.T) {

	t.Parallel()

	var d db.Database = &Service{}
	if d == nil {
		t.Fatal("Does not implement db.Database interface.")
	}

}

func TestInterfaceService(t *testing.T) {

	t.Parallel()

	var s m.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}

func openDatabase(t *testing.T) *Service {

	f, err := ioutil.TempFile(empty, "bolt-")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s", err.Error())
	}
	filename := f.Name()
	f.Close()

	database := NewService(&db.Configuration{
		Location: filename,
	})
	err = database.Start()
	assertNoError(t, err)
	if !database.running {
		t.Error("database is not running")
	}

	return database

}

func closeDatabase(t *testing.T, database *Service) {

	// close database
	database.Stop()
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

func closeDatabaseX(t *testing.T, database *Service) {

	// close database
	database.Stop()
	if database.running {
		t.Error("database is still running")
	}
	if database.db != nil {
		t.Error("database.db is not nil")
	}

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
