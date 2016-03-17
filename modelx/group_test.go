package model

import (
	"testing"
)

func TestGroup(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	fakeUsername := "jeff@lubowski.zzz"
	fakePassword := "12345678"
	fakeUserID := ""
	fakeGroupname1 := "Group1"
	fakeGroupname2 := "Group2"

	// add user to database
	if err := database.Update(func(tx Transaction) error {
		user := NewUser(fakeUsername)
		if err := user.SetPassword(fakePassword); err != nil {
			return err
		}
		if err := user.Save(tx); err != nil {
			return err
		}
		fakeUserID = user.ID
		return nil
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// add group1
	if err := database.Update(func(tx Transaction) error {
		group := NewGroup(fakeUserID, fakeGroupname1)
		return group.Save(tx)
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve group1 by name
	if err := database.Select(func(tx Transaction) error {
		group, err := GroupByName(fakeUserID, fakeGroupname1, tx)
		if err != nil {
			return err
		}
		if group == nil {
			t.Errorf("Group1 not found by name: %s", fakeGroupname1)
		} else if group.Name != fakeGroupname1 {
			t.Errorf("Bad name, expected %s, actual %s", fakeGroupname1, group.Name)
		} else if group.ID == empty {
			t.Error("Empty Group ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// retrieve groups by user
	if err := database.Select(func(tx Transaction) error {
		groups, err := GroupsByUser(fakeUserID, tx)
		if err != nil {
			return err
		}
		if groups == nil {
			t.Error("Groups not found by user")
		} else if len(groups) != 1 {
			t.Errorf("Bad number of groups, expected %d, actual %d", 1, len(groups))
		} else if groups.First().Name != fakeGroupname1 {
			t.Errorf("Bad name, expected %s, actual %s", fakeGroupname1, groups.First().Name)
		} else if groups.First().ID == empty {
			t.Error("Empty Group ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// add group2
	if err := database.Update(func(tx Transaction) error {
		group := NewGroup(fakeUserID, fakeGroupname2)
		return group.Save(tx)
	}); err != nil {
		t.Fatalf("Error updating database: %s", err.Error())
	}

	// retrieve group2 by name
	if err := database.Select(func(tx Transaction) error {
		group, err := GroupByName(fakeUserID, fakeGroupname2, tx)
		if err != nil {
			return err
		}
		if group == nil {
			t.Errorf("Group2 not found by name: %s", fakeGroupname2)
		} else if group.Name != fakeGroupname2 {
			t.Errorf("Bad name, expected %s, actual %s", fakeGroupname2, group.Name)
		} else if group.ID == empty {
			t.Error("Empty Group ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// retrieve groups by user
	if err := database.Select(func(tx Transaction) error {
		groups, err := GroupsByUser(fakeUserID, tx)
		if err != nil {
			return err
		}
		if groups == nil {
			t.Error("Groups not found by user")
		} else if len(groups) != 2 {
			t.Errorf("Bad number of groups, expected %d, actual %d", 2, len(groups))
		} else if groups.First().Name != fakeGroupname1 {
			t.Errorf("Bad name, expected %s, actual %s", fakeGroupname1, groups.First().Name)
		} else if groups.First().ID == empty {
			t.Error("Empty Group ID")
		} else if groups[1].Name != fakeGroupname2 {
			t.Errorf("Bad name, expected %s, actual %s", fakeGroupname2, groups[1].Name)
		} else if groups[1].ID == empty {
			t.Error("Empty Group ID")
		}
		return nil
	}); err != nil {
		t.Fatalf("Error selecting database: %s", err.Error())
	}

	// add new user with same username to database, expect error
	if err := database.Update(func(tx Transaction) error {
		group := NewGroup(fakeUserID, fakeGroupname1)
		return group.Save(tx)
	}); err == nil {
		t.Error("Expected error saving group with non-unique name.")
	}

}
