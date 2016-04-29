package model

import (
	"fmt"
	"testing"
)

func TestGroupSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityGroup); obj == nil {
		t.Error("missing getObject entry")
	} else if !obj.hasIncrementingID() {
		t.Error("groups have incrementing IDs")
	}

	if obj := allEntities[entityGroup]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestGroupIncrementingID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	if err := db.Update(func(tx Transaction) error {
		g := G.New("0000000001", "G1")
		return G.Save(tx, g)
	}); err != nil {
		t.Fatalf("error adding new group: %s", err.Error())
	}

	if err := db.Update(func(tx Transaction) error {
		id, err := tx.NextID(entityGroup)
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

func TestGroupGetBadID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if group := G.Get(tx, empty); group != nil {
			t.Error("Expected nil group")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestGroups(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	userID := empty

	err := db.Update(func(tx Transaction) error {

		user := U.New("user1", "password")
		if err := U.Save(tx, user); err != nil {
			return err
		}
		userID = user.ID

		for i := 0; i < 20; i++ {
			g := G.New(userID, fmt.Sprintf("Group%d", i))
			if err := G.Save(tx, g); err != nil {
				return err
			}
		}

		return nil

	})
	if err != nil {
		t.Fatalf("Error adding groups: %s", err.Error())
	}

	err = db.Update(func(tx Transaction) error {

		groups := G.GetForUser(tx, userID)
		if groups == nil {
			t.Fatal("Nil group, expected non-nil value")
		}

		if len(groups) != 20 {
			t.Errorf("bad number of groups, expected %d, actual %d", 20, len(groups))
		}

		for i, group := range groups {
			if i%2 == 0 {
				if err := G.Delete(tx, group.ID); err != nil {
					return err
				}
			}
		}

		return nil

	})
	if err != nil {
		t.Fatalf("Error deleting groups: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {

		groups := G.GetForUser(tx, userID)
		if groups == nil {
			t.Fatal("Nil group, expected non-nil value")
		}

		if len(groups) != 10 {
			t.Errorf("bad number of groups, expected %d, actual %d", 20, len(groups))
		}

		return nil

	})
	if err != nil {
		t.Fatalf("Error reading groups: %s", err.Error())
	}

}
