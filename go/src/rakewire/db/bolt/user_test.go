package bolt

import (
	m "rakewire/model"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	username := "jeff@jarvis.com"
	password := "abcdefg"

	u := m.NewUser(username)
	u.SetPassword(password)

	if u.ID != 0 {
		t.Errorf("New users must have an ID of 0, actual value: %d", u.ID)
	}

	if err := db.UserSave(u); err != nil {
		t.Fatalf("Error saving user: %s", err.Error())
	}

	if u.ID == 0 {
		t.Error("UserID not set after save")
	}

	if !u.MatchPassword(password) {
		t.Error("user password does not match")
	}

	if u2, err := db.UserGetByUsername(username); err != nil {
		t.Errorf("Error retrieving user by username: %s", err.Error())
	} else if u2 == nil {
		t.Errorf("User by username not found: %s", username)
	} else {
		if u2.ID != u.ID {
			t.Errorf("user IDs do not match, expected: %s, actial: %s", u.ID, u2.ID)
		}
		if u2.Username != u.Username {
			t.Errorf("usernames not not match, expected: %s, actial: %s", u.Username, u2.Username)
		}
		if !u2.MatchPassword(password) {
			t.Error("user password does not match")
		}
	}

	if u2, err := db.UserGetByUsername(strings.ToUpper(username)); err != nil {
		t.Errorf("Retrieving user by username is not case-insensitive: %s", err.Error())
	} else if u2 == nil {
		t.Errorf("User by username not found: %s", username)
	} else {
		if u2.ID != u.ID {
			t.Errorf("user IDs do not match, expected: %s, actial: %s", u.ID, u2.ID)
		}
		if u2.Username != u.Username {
			t.Errorf("usernames not not match, expected: %s, actial: %s", u.Username, u2.Username)
		}
		if !u2.MatchPassword(password) {
			t.Error("user password does not match")
		}
	}

	if u2, err := db.UserGetByFeverHash(u.FeverHash); err != nil {
		t.Errorf("Error retrieving user by feverhash: %s", err.Error())
	} else if u2 == nil {
		t.Errorf("User by feverhash not found: %s", u.FeverHash)
	} else {
		if u2.ID != u.ID {
			t.Errorf("user IDs not not match, expected: %s, actial: %s", u.ID, u2.ID)
		}
		if u2.Username != u.Username {
			t.Errorf("usernames do not match, expected: %s, actial: %s", u.Username, u2.Username)
		}
		if !u2.MatchPassword(password) {
			t.Error("user password does not match")
		}
	}

}

func TestUserUniqueUsername(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	username := "Jeff@Jarvis.com"
	password := "abcdefg"

	u := m.NewUser(username)
	u.SetPassword(password)
	if err := db.UserSave(u); err != nil {
		t.Fatalf("Error saving user: %s", err.Error())
	}

	u2 := m.NewUser(strings.ToUpper(username))
	u2.SetPassword(password)
	err := db.UserSave(u2)
	if err == nil {
		t.Error("Expected error, none returned")
	} else if err.Error() != "Cannot save user, username is already taken: "+strings.ToLower(username) {
		t.Errorf("Bad error text: %s", err.Error())
	}

}
