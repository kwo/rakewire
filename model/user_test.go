package model

import (
	"strings"
	"testing"
)

func TestUserSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityUser); obj == nil {
		t.Error("missing getObject entry")
	} else if !obj.hasIncrementingID() {
		t.Error("users have incrementing IDs")
	}

	if obj := allEntities[entityUser]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestUserIncrementingID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	if err := db.Update(func(tx Transaction) error {
		u := U.New("username", "hello")
		return U.Save(tx, u)
	}); err != nil {
		t.Fatalf("error adding new user: %s", err.Error())
	}

	if err := db.Update(func(tx Transaction) error {
		id, err := tx.NextID(entityUser)
		if err != nil {
			return err
		}
		if id != 2 {
			t.Errorf("bad next id: %d, expected %d", id, 2)
		}
		return nil
	}); err != nil {
		t.Fatalf("error getting next id: %s", err.Error())
	}

}

func TestUserPassword(t *testing.T) {

	t.Parallel()

	username := "jeff@jarvis.com"
	password := "abcedfghijklmnopqrstuvwxyz"

	u := U.New(username, password)
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

func TestUserGetBadID(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	err := database.Select(func(tx Transaction) error {
		if user := U.Get(tx, empty); user != nil {
			t.Error("Expected nil user")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestUserGetBadUsername(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	err := database.Select(func(tx Transaction) error {
		if user := U.GetByUsername(tx, empty); user != nil {
			t.Error("Expected nil user")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestUserGetBadFeverhash(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	err := database.Select(func(tx Transaction) error {
		if user := U.GetByFeverhash(tx, empty); user != nil {
			t.Error("Expected nil user")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
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
		user := U.New(fakeUsername, fakePassword)
		if err := U.Save(tx, user); err != nil {
			return err
		}
		fakeID = user.ID
		return nil
	}); err != nil {
		t.Fatalf("Error adding user: %s", err.Error())
	}

	// retrieve user by ID
	if err := database.Select(func(tx Transaction) error {
		user := U.Get(tx, fakeID)
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
		if err := U.Delete(tx, fakeID); err != nil {
			t.Errorf("Error deleting user: %s", err.Error())
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve user by ID
	if err := database.Select(func(tx Transaction) error {
		user := U.Get(tx, fakeID)
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
		user := U.New(fakeUsername, fakePassword)
		fakeFeverhash = user.FeverHash
		return U.Save(tx, user)
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve user by username
	if err := database.Select(func(tx Transaction) error {
		user := U.GetByUsername(tx, fakeUsername)
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
		user := U.GetByFeverhash(tx, fakeFeverhash)
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
		user := U.New(strings.ToUpper(fakeUsername), fakePassword) // test case-insensitivity
		if err := U.Save(tx, user); err != ErrUsernameTaken {
			t.Error("Expected error saving user with non-unique username.")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

}

func TestUserJson(t *testing.T) {

	t.Parallel()

	user := U.New("karl@ostendorf.com", "abcdefg")
	data, err := user.encode()
	if err != nil {
		t.Errorf("Error encoding user: %s", err.Error())
	}

	t.Log(string(data))

}
