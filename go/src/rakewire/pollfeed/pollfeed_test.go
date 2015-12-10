package pollfeed

import (
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	m "rakewire/model"
	"testing"
	"time"
)

const (
	databaseFile = "../../../test/pollfeed.db"
)

func TestInterfaceService(t *testing.T) {

	var s m.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}

func TestPoll(t *testing.T) {

	//t.SkipNow()

	// open database
	database := bolt.NewService(&db.Configuration{
		Location: databaseFile,
	})
	err := database.Start()
	assertNoError(t, err)

	// create service
	cfg := &Configuration{
		Interval: "1m",
	}
	pf := NewService(cfg, database)
	pf.pollInterval = 50 * time.Millisecond

	pf.Start()
	assertEqual(t, true, pf.IsRunning())
	time.Sleep(100 * time.Millisecond)
	pf.Stop()
	assertEqual(t, false, pf.IsRunning())

	// close database
	database.Stop()

	// remove file
	err = os.Remove(databaseFile)
	assertNoError(t, err)

}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Not equal: expected %v, actual %v", a, b)
	}
}
