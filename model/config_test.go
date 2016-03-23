package model

import (
	"testing"
)

func TestConfigSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityConfig); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityConfig]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestConfig(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	loglevel := "TRACE"
	var userID uint64 = 1

	// get, update config
	if err := database.Update(func(tx Transaction) error {
		config := C.Get(tx)
		config.Log.Level = loglevel
		config.Sequences.User = userID
		return C.Put(tx, config)
	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

	// get, compare config
	if err := database.Select(func(tx Transaction) error {
		config := C.Get(tx)

		if config.Log.Level != loglevel {
			t.Errorf("Bad loglevel, expected %s, actual %s", loglevel, config.Log.Level)
		}

		if config.Sequences.User != userID {
			t.Errorf("Bad user sequence, expected %d, actual %d", config.Sequences.User, userID)
		}

		return nil

	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

}
