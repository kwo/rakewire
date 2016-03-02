package store

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestBoltDB(t *testing.T) {

	t.Parallel()

	store := openTestStore(t)
	defer closeTestStore(t, store)

}

func openTestStore(t *testing.T, flags ...bool) Store {

	//flagPopulateStore := len(flags) > 0 && flags[0]

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	store, err := Instance.Open(location)
	if err != nil {
		t.Fatalf("Cannot open store: %s", err.Error())
	}

	// if flagPopulateStore {
	// 	err = boltDB.Update(func(tx Transaction) error {
	// 		return populateStore(t, tx)
	// 	})
	// 	if err != nil {
	// 		t.Fatalf("Cannot populate store: %s", err.Error())
	// 	}
	// }

	return store

}

func closeTestStore(t *testing.T, store Store) {

	location := store.Location()

	if err := Instance.Close(store); err != nil {
		t.Errorf("Cannot close store: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}
