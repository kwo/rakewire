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

	// add test users
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

	// add groups for both users
	for i := 0; i < 10; i++ {

		group1 := m.NewGroup(user1.ID, fmt.Sprintf("User%d-Group%d", 1, i))
		if err := db.GroupSave(group1); err != nil {
			t.Fatalf("Error saving group1: %s", err.Error())
		}
		if g, err := db.GroupGet(group1.ID); err != nil {
			t.Errorf("Error retrieving group1: %s", err.Error())
		} else if g == nil {
			t.Errorf("Group not added: %d", group1.ID)
		}

		group2 := m.NewGroup(user2.ID, fmt.Sprintf("User%d-Group%d", 2, i))
		if err := db.GroupSave(group2); err != nil {
			t.Fatalf("Error saving group2: %s", err.Error())
		}
		if g, err := db.GroupGet(group2.ID); err != nil {
			t.Errorf("Error retrieving group2: %s", err.Error())
		} else if g == nil {
			t.Errorf("Group not added: %d", group2.ID)
		}

	}

	// check that groups are there
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

	// delete odd-numbered groups
	for i := 0; i < 10; i++ {

		if i%2 == 0 {
			continue
		}

		if err := db.GroupDelete(groups1[i]); err != nil {
			t.Errorf("Cannot delete from group1: %s", err.Error())
		}
		if g, err := db.GroupGet(groups1[i].ID); err != nil {
			t.Errorf("Error retrieving group1: %s", err.Error())
		} else if g != nil {
			t.Errorf("Group not deleted: %d", groups1[i].ID)
		}

		if err := db.GroupDelete(groups2[i]); err != nil {
			t.Errorf("Cannot delete from group2: %s", err.Error())
		}
		if g, err := db.GroupGet(groups2[i].ID); err != nil {
			t.Errorf("Error retrieving group2: %s", err.Error())
		} else if g != nil {
			t.Errorf("Group not deleted: %d", groups2[i].ID)
		}

	}

	// refetch groups for each test user
	groups1, err = db.GroupGetAllByUser(user1.ID)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err.Error())
	}
	if len(groups1) != 5 {
		t.Fatalf("Group count mismatch, expected %d, actual %d", 5, len(groups1))
	}

	groups2, err = db.GroupGetAllByUser(user2.ID)
	if err != nil {
		t.Fatalf("Failed to get groups: %s", err.Error())
	}
	if len(groups2) != 5 {
		t.Fatalf("Group count mismatch, expected %d, actual %d", 5, len(groups2))
	}

	// check that odd-numbered groups are deleted
	for i := 0; i < 5; i++ {

		name1 := fmt.Sprintf("User%d-Group%d", 1, i*2)
		if groups1[i].Name != name1 {
			t.Errorf("Bad group name, expected %s, actual %s", name1, groups1[i].Name)
		}

		name2 := fmt.Sprintf("User%d-Group%d", 2, i*2)
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
