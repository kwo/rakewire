package modelng

import (
	"testing"
)

func TestConfig(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	loglevel := "TRACE"
	var userID uint64 = 1

	// get, update config
	if err := database.Update(func(tx Transaction) error {
		config := C.Get(tx)
		config.LoggingLevel = loglevel
		config.Sequences.User = userID
		return C.Put(config, tx)
	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

	// get, compare config
	if err := database.Select(func(tx Transaction) error {
		config := C.Get(tx)

		if config.LoggingLevel != loglevel {
			t.Errorf("Bad loglevel, expected %s, actual %s", loglevel, config.LoggingLevel)
		}

		if config.Sequences.User != userID {
			t.Errorf("Bad user sequence, expected %d, actual %d", config.Sequences.User, userID)
		}

		return nil

	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

}
