package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestRenameWithTimestamp(t *testing.T) {

	t.Parallel()

	tmpdir := os.TempDir()
	filename := "data.db"
	location := filepath.Join(tmpdir, filename)
	if err := ioutil.WriteFile(location, []byte{}, os.ModePerm); err != nil {
		t.Fatalf("Error touching file: %s", err.Error())
	}

	if newLocation, err := renameWithTimestamp(location); err != nil {
		t.Errorf("Error renaming file: %s", err.Error())
	} else {

		t.Log(newLocation)

		if _, err := os.Stat(location); err == nil {
			t.Error("Expected error getting stats for now non-existent original file")
		}

		if _, err := os.Stat(newLocation); err != nil {
			t.Errorf("Error getting stats for new file: %s", err.Error())
		}

		if err := os.Remove(newLocation); err != nil {
			t.Errorf("Cannot remove temp file: %s", err.Error())
		}

	}

}

func TestContainerIterate(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	if err := database.Update(func(tx Transaction) error {
		feed := NewFeed("http://localhost/")
		_, err := feed.Save(tx)
		return err
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	if err := database.Select(func(tx Transaction) error {
		feeds, err := tx.Container("Data/Feed")
		if err != nil {
			return err
		}
		return feeds.Iterate(func(record Record) error {
			for k, v := range record {
				t.Logf("%s: %s", k, v)
			}
			return nil
		})
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

}
