package bolt

import (
	"fmt"
	m "rakewire/model"
	"testing"
)

func TestGroup(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	user1 := m.NewUser("jeff@jarvis.com")
	user1.SetPassword("abcdefg")
	if err := db.UserSave(user1); err != nil {
		t.Fatalf("Error saving user1: %s", err.Error())
	}

	user2 := m.NewUser("nick@taleb.com")
	user2.SetPassword("abcdefg")
	if err := db.UserSave(user2); err != nil {
		t.Fatalf("Error saving user2: %s", err.Error())
	}

	for i := 0; i < 10; i++ {

		group1 := m.NewGroup(user1.ID, fmt.Sprintf("User%d-Group%d", 1, i))
		if err := db.GroupSave(group1); err != nil {
			t.Fatalf("Error saving group1: %s", err.Error())
		}
		if group1.ID == 0 {
			t.Error("New Group.ID not set on save")
		}

		group2 := m.NewGroup(user2.ID, fmt.Sprintf("User%d-Group%d", 2, i))
		if err := db.GroupSave(group2); err != nil {
			t.Fatalf("Error saving group2: %s", err.Error())
		}
		if group2.ID == 0 {
			t.Error("New Group.ID not set on save")
		}

	}

	groups1, err := db.GroupGetAllByUser(user1.ID)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err.Error())
	}
	if len(groups1) != 10 {
		t.Fatalf("Group count mismatch, expected %d, actual %d", 10, len(groups1))
	}

	groups2, err := db.GroupGetAllByUser(user2.ID)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err.Error())
	}
	if len(groups2) != 10 {
		t.Fatalf("Group count mismatch, expected %d, actual %d", 10, len(groups2))
	}

	for i := 0; i < 10; i++ {

		name1 := fmt.Sprintf("User%d-Group%d", 1, i)
		if groups1[i].Name != name1 {
			t.Errorf("Bad group name, expected %s, actual %s", name1, groups1[i].Name)
		}

		name2 := fmt.Sprintf("User%d-Group%d", 2, i)
		if groups2[i].Name != name2 {
			t.Errorf("Bad group name, expected %s, actual %s", name2, groups2[i].Name)
		}

	}

}

func TestGroupUniqueName(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	user := m.NewUser("Jeff@Jarvis.com")
	user.SetPassword("abcdefg")
	if err := db.UserSave(user); err != nil {
		t.Fatalf("Error saving user: %s", err.Error())
	}

	g1 := m.NewGroup(user.ID, "name1")
	if err := db.GroupSave(g1); err != nil {
		t.Fatalf("Error saving group: %s", err.Error())
	}

	g2 := m.NewGroup(user.ID, "name1")
	err := db.GroupSave(g2)
	if err == nil {
		t.Error("Expected error, none returned")
	} else if err.Error() != "Cannot save group, group name is already taken: name1" {
		t.Errorf("Bad error text: %s", err.Error())
	}

}
