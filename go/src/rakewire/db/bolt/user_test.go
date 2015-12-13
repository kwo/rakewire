package bolt

import (
	m "rakewire/model"
	"testing"
)

func TestUser(t *testing.T) {

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	username := "jeff@jarvis.com"
	password := "abcdefg"

	u := m.NewUser(username)
	u.SetPassword(password)
	if err := db.UserSave(u); err != nil {
		t.Fatalf("Error saving user: %s", err.Error())
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
