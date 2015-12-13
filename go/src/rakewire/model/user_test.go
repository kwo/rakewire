package model

import (
	"testing"
)

func TestUserPassword(t *testing.T) {

	username := "jeff@jarvis.com"
	password := "abcedfghijklmnopqrstuvwxyz"

	u := NewUser(username)
	if len(u.ID) != 36 {
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
