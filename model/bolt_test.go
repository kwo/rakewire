package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestBoltDB(t *testing.T) {

	t.Parallel()

	store := openTestDatabase(t)
	defer closeTestDatabase(t, store)

}

func openTestDatabase(t *testing.T, flags ...bool) Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	store, err := Instance.Open(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	return store

}

func closeTestDatabase(t *testing.T, db Database) {

	location := db.Location()

	if err := Instance.Close(db); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func TestWhatAreUniqueKeys(t *testing.T) {

	t.Parallel()

	var uniques []string

	// collect unique indexes
	for entityName := range allEntities {
		o := getObject(entityName)
		for indexName, values := range o.indexes() {
			if len(values) < 3 {
				uniques = append(uniques, fmt.Sprintf("%s:%s", entityName, indexName))
			}
		} // indexes
	} // entities

	for _, indexName := range uniques {
		t.Log(indexName)
	}

}
