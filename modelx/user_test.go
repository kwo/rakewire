package model

import (
	"testing"
)

func TestUserPassword(t *testing.T) {

	username := "jeff@jarvis.com"
	password := "abcedfghijklmnopqrstuvwxyz"

	u := NewUser(username)
	if u.ID != empty {
		t.Error("User ID not set properly by factory method")
	}
	if u.Username != username {
		t.Error("Username not set properly by factory method")
	}

	if err := u.SetPassword(password); err != nil {
		t.Errorf("Cannot set password: %s", err.Error())
	}
	if !u.MatchPassword(password) {
		t.Error("User passwords should match")
	}
	if u.MatchPassword("abcde") {
		t.Error("User passwords should NOT match")
	}

}

func TestUserIndexes(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	fakeUsername := "jeff@lubowski.zzz"
	fakePassword := "12345678"
	fakeFeverhash := ""

	// add user to database
	if err := database.Update(func(tx Transaction) error {
		user := NewUser(fakeUsername)
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		fakeFeverhash = user.FeverHash
		return user.Save(tx)
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve user by username
	if err := database.Select(func(tx Transaction) error {
		user, err := UserByUsername(fakeUsername, tx)
		if err != nil {
			return err
		}
		if user == nil {
			t.Errorf("User not found by username: %s", fakeUsername)
		} else if user.Username != fakeUsername {
			t.Errorf("Bad username, expected %s, actual %s", fakeUsername, user.Username)
		} else if user.ID == empty {
			t.Error("Empty User ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// retrieve user by feverhash
	if err := database.Select(func(tx Transaction) error {
		user, err := UserByFeverHash(fakeFeverhash, tx)
		if err != nil {
			return err
		}
		if user == nil {
			t.Errorf("User not found by feverhash: %s", fakeUsername)
		} else if user.Username != fakeUsername {
			t.Errorf("Bad username, expected %s, actual %s", fakeUsername, user.Username)
		} else if user.ID == empty {
			t.Error("Empty User ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// add new user with same username to database, expect error
	if err := database.Update(func(tx Transaction) error {
		user := NewUser(fakeUsername)
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		if err := user.Save(tx); err == nil {
			t.Error("Expected error saving user with non-unique username.")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

}
