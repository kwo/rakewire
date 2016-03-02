package modelng

import (
	"strings"
	"testing"
)

func TestUserPassword(t *testing.T) {

	t.Parallel()

	username := "jeff@jarvis.com"
	password := "abcedfghijklmnopqrstuvwxyz"

	u := U.New(username)
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

func TestUserAddDelete(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	fakeID := ""
	fakeUsername := "jeff@lubowski.zzz"
	fakePassword := "12345678"

	// add user to database
	if err := database.Update(func(tx Transaction) error {
		user := U.New(fakeUsername)
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		if err := U.Save(user, tx); err != nil {
			return err
		}
		fakeID = user.ID
		return nil
	}); err != nil {
		t.Fatalf("Error adding user: %s", err.Error())
	}

	// retrieve user by ID
	if err := database.Select(func(tx Transaction) error {
		user := U.GetByID(fakeID, tx)
		if user == nil {
			t.Errorf("User not found by ID: %s", fakeID)
		} else if user.Username != fakeUsername {
			t.Errorf("Bad username, expected %s, actual %s", fakeUsername, user.Username)
		} else if user.ID != fakeID {
			t.Errorf("Bad ID, expected %s, actual %s", fakeID, user.ID)
		}
		return nil
	}); err != nil {
		t.Fatalf("Error retrieving user by ID: %s", err.Error())
	}

	// delete user
	if err := database.Update(func(tx Transaction) error {
		if err := U.Delete(fakeID, tx); err != nil {
			t.Errorf("Error deleting user: %s", err.Error())
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve user by ID
	if err := database.Select(func(tx Transaction) error {
		user := U.GetByID(fakeID, tx)
		if user != nil {
			t.Error("User found by ID, expected nil")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error retrieving user by ID: %s", err.Error())
	}

}

func TestUserGets(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	fakeUsername := "jeff@lubowski.zzz"
	fakePassword := "12345678"
	fakeFeverhash := ""

	// add user to database
	if err := database.Update(func(tx Transaction) error {
		user := U.New(fakeUsername)
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		fakeFeverhash = user.FeverHash
		return U.Save(user, tx)
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve user by username
	if err := database.Select(func(tx Transaction) error {
		user := U.GetByUsername(fakeUsername, tx)
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
		user := U.GetByFeverhash(fakeFeverhash, tx)
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
		user := U.New(strings.ToUpper(fakeUsername)) // test case-insensitivity
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		if err := U.Save(user, tx); err != ErrUsernameTaken {
			t.Error("Expected error saving user with non-unique username.")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

}
