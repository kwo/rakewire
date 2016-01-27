package reaper

import (
	"io/ioutil"
	"os"
	"rakewire/model"
	"testing"
	"time"
)

const (
	databaseFile = "../test/pollfeed.db"
)

func TestInterfaceService(t *testing.T) {

	var s model.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}

func TestReaper(t *testing.T) {

	t.SkipNow()

	// open database
	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	// create service
	cfg := model.NewConfiguration()
	r := NewService(cfg, database)

	r.Start()
	if !r.IsRunning() {
		t.Error("Polling service is not running")
	}
	time.Sleep(100 * time.Millisecond)
	r.Stop()
	if r.IsRunning() {
		t.Error("Polling service is still running")
	}

}

func openTestDatabase(t *testing.T) model.Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	boltDB, err := model.OpenDatabase(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d model.Database) {

	location := d.Location()

	if err := model.CloseDatabase(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}
