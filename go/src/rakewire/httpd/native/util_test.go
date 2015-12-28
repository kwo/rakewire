package native

import (
	"io/ioutil"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	"rakewire/logging"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := &logging.Configuration{Level: logging.LogWarn}
	cfg.Init()
	status := m.Run()
	os.Exit(status)
}

func openDatabase(t *testing.T) (*bolt.Service, string) {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s\n", err.Error())
	}
	testDatabaseFile := f.Name()
	f.Close()

	cfg := db.Configuration{
		Location: testDatabaseFile,
	}
	testDatabase := bolt.NewService(&cfg)
	err = testDatabase.Start()
	if err != nil {
		t.Fatalf("Cannot open database: %s\n", err.Error())
	}

	return testDatabase, testDatabaseFile

}

func closeDatabase(t *testing.T, database *bolt.Service, testDatabaseFile string) {

	database.Stop()
	if err := os.Remove(testDatabaseFile); err != nil {
		t.Errorf("Cannot delete temp database file: %s", err.Error())
	}

}
