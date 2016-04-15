package rest

import (
	"io/ioutil"
	"os"
	"rakewire/model"
	"testing"
)

func TestMain(m *testing.M) {
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
