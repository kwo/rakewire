package model

import (
	"testing"
)

func TestSubscriptions(t *testing.T) {

	subscription := NewSubscription(1, 1)

	subscription.AddGroup(3)
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.AddGroup(2)
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	subscription.AddGroup(2)
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	if found := subscription.HasGroup(2); !found {
		t.Error("Group not found: 2")
	}

	if found := subscription.HasGroup(3); !found {
		t.Error("Group not found: 3")
	}

	if found := subscription.HasGroup(1); found {
		t.Error("Unexpected group found: 1")
	}

	subscription.RemoveGroup(3)
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.RemoveGroup(2)
	if len(subscription.GroupIDs) != 0 {
		t.Errorf("Bad group count: expcected %d, actual %d", 0, len(subscription.GroupIDs))
	}

}
