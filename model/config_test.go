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

	var userID uint64 = 1

	// get, update config
	if err := database.Update(func(tx Transaction) error {
		config := C.Get(tx)
		config.SetBool("one", true)
		config.SetInt("two", 1)
		config.SetStr("three", empty)
		config.Sequences.User = userID
		return C.Put(tx, config)
	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

	// get, compare config
	if err := database.Select(func(tx Transaction) error {
		config := C.Get(tx)

		boolValue := config.GetBool("one")
		boolValueExpected := true
		if !boolValue {
			t.Errorf("Bad bool value: %t, expected %t", boolValue, boolValueExpected)
		}

		intValue := config.GetInt("two")
		intValueExpected := 1
		if intValue != intValueExpected {
			t.Errorf("Bad int value: %d, expected %d", intValue, intValueExpected)
		}

		strValue := config.GetStr("three", "hello")
		strValueExpected := "hello"
		if strValue != strValueExpected {
			t.Errorf("Bad str value: %s, expected %s", strValue, strValueExpected)
		}

		if config.Sequences.User != userID {
			t.Errorf("Bad user sequence, expected %d, actual %d", config.Sequences.User, userID)
		}

		return nil

	}); err != nil {
		t.Fatalf("Error retrieving config: %s", err.Error())
	}

}

func TestConfigJson(t *testing.T) {

	t.Parallel()
	config := C.New()
	config.Sequences.User++
	config.SetStr("logging.level", "DEBUG")

	data, err := config.encode()
	if err != nil {
		t.Errorf("Error encoding config: %s", err.Error())
	}

	t.Log(string(data))

}
