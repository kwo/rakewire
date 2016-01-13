package model

import (
	"github.com/boltdb/bolt"
	"testing"
)

func TestConfig(t *testing.T) {

	db, err := openDatabase()
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	defer closeDatabase(db)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(configurationBucketName))
		return err
	})
	if err != nil {
		t.Fatalf("Cannot prepare database: %s", err.Error())
	}

	database := NewBoltDatabase(db)

	config := NewConfiguration()
	config.Set("hello", "world")
	err = database.Update(func(tx Transaction) error {
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

	if value := config.Get("hello"); value != "world" {
		t.Errorf("Expected %s, actual %s", "world", value)
	}

}
