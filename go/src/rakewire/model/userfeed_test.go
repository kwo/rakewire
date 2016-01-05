package model

import (
	"testing"
)

func TestUserFeeds(t *testing.T) {

	userfeed := NewUserFeed(1, 1)

	userfeed.AddGroup(3)
	if len(userfeed.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(userfeed.GroupIDs))
	}

	userfeed.AddGroup(2)
	if len(userfeed.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(userfeed.GroupIDs))
	}

	userfeed.AddGroup(2)
	if len(userfeed.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(userfeed.GroupIDs))
	}

	if found := userfeed.HasGroup(2); !found {
		t.Error("Group not found: 2")
	}

	if found := userfeed.HasGroup(3); !found {
		t.Error("Group not found: 3")
	}

	if found := userfeed.HasGroup(1); found {
		t.Error("Unexpected group found: 1")
	}

	userfeed.RemoveGroup(3)
	if len(userfeed.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(userfeed.GroupIDs))
	}

	userfeed.RemoveGroup(2)
	if len(userfeed.GroupIDs) != 0 {
		t.Errorf("Bad group count: expcected %d, actual %d", 0, len(userfeed.GroupIDs))
	}

}
