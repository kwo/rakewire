package rest

import (
	"io/ioutil"
	"os"
	"rakewire/logging"
	"rakewire/model"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := &logging.Configuration{Level: logging.LogWarn}
	cfg.Init()
	status := m.Run()
	os.Exit(status)
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

	// err = boltDB.Update(func(tx model.Transaction) error {
	// 	return populateDatabase(tx)
	// })
	// if err != nil {
	// 	t.Fatalf("Cannot populate database: %s", err.Error())
	// }

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