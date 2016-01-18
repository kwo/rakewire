package model

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestBoltDB(t *testing.T) {

	db := openTestDatabase(t)
	closeTestDatabase(t, db)

}

func openTestDatabase(t *testing.T) Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	boltDB, err := OpenDatabase(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d Database) {

	location := d.Location()

	if err := CloseDatabase(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}
