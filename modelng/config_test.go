package modelng

import (
	"testing"
)

func TestConfig(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	fakeName := "sequence"
	fakeValue := "001"

	// add config entry to database
	if err := database.Update(func(tx Transaction) error {
		entry := C.New(fakeName, fakeValue)
		if err := C.Save(entry, tx); err != nil {
			t.Fatalf("Error saving config entry: %s", err.Error())
		}
		return nil
	}); err != nil {
		t.Fatalf("Error adding config entry: %s", err.Error())
	}

	// retrieve entry by ID
	if err := database.Select(func(tx Transaction) error {
		entry := C.GetByID(fakeName, tx)
		if entry == nil {
			t.Errorf("Config entry not found by ID: %s", fakeName)
		} else if entry.Value != fakeValue {
			t.Errorf("Bad value, expected %s, actual %s", fakeValue, entry.Value)
		} else if entry.Name != fakeName {
			t.Errorf("Bad name, expected %s, actual %s", fakeName, entry.Name)
		}
		return nil
	}); err != nil {
		t.Fatalf("Error retrieving entry by ID: %s", err.Error())
	}

	// delete entry
	if err := database.Update(func(tx Transaction) error {
		if err := C.Delete(fakeName, tx); err != nil {
			t.Errorf("Error deleting entry: %s", err.Error())
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve entry by ID
	if err := database.Select(func(tx Transaction) error {
		entry := U.GetByID(fakeName, tx)
		if entry != nil {
			t.Error("Entry found by ID, expected nil")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error retrieving entry by ID: %s", err.Error())
	}

}
