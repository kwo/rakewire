package pollfeed

import (
	"github.com/kwo/rakewire/model"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestInterfaceService(t *testing.T) {

	var s model.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement model.Service interface.")
	}

}

func TestPoll(t *testing.T) {

	// open database
	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	// create service
	cfg := &Configuration{
		BatchMax:        10,
		IntervalSeconds: 1,
	}
	pf := NewService(cfg, database)

	pf.Start()
	if !pf.IsRunning() {
		t.Error("Polling service is not running")
	}
	time.Sleep(100 * time.Millisecond)
	pf.Stop()
	if pf.IsRunning() {
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

	boltDB, err := model.Instance.Open(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d model.Database) {

	location := d.Location()

	if err := model.Instance.Close(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}
