package model

import (
	"testing"
)

func TestConfig(t *testing.T) {

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	config := NewConfiguration()
	config.Set("hello", "world")
	err := database.Update(func(tx Transaction) error {
		if err := config.Save(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot save to database: %s", err.Error())
	}

	config = NewConfiguration()
	err = database.Select(func(tx Transaction) error {
		if err := config.Load(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot load from database: %s", err.Error())
	}

	if value := config.Get("hello", ""); value != "world" {
		t.Errorf("Expected %s, actual %s", "world", value)
	}

}
