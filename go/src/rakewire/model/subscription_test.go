package model

import (
	"testing"
)

func TestSubscriptions(t *testing.T) {

	subscription := NewSubscription(kvKeyUintEncode(1), kvKeyUintEncode(1))

	subscription.AddGroup(kvKeyUintEncode(3))
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.AddGroup(kvKeyUintEncode(2))
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	subscription.AddGroup(kvKeyUintEncode(2))
	if len(subscription.GroupIDs) != 2 {
		t.Errorf("Bad group count: expcected %d, actual %d", 2, len(subscription.GroupIDs))
	}

	if found := subscription.HasGroup(kvKeyUintEncode(2)); !found {
		t.Error("Group not found: 2")
	}

	if found := subscription.HasGroup(kvKeyUintEncode(3)); !found {
		t.Error("Group not found: 3")
	}

	if found := subscription.HasGroup(kvKeyUintEncode(1)); found {
		t.Error("Unexpected group found: 1")
	}

	subscription.RemoveGroup(kvKeyUintEncode(3))
	if len(subscription.GroupIDs) != 1 {
		t.Errorf("Bad group count: expcected %d, actual %d", 1, len(subscription.GroupIDs))
	}

	subscription.RemoveGroup(kvKeyUintEncode(2))
	if len(subscription.GroupIDs) != 0 {
		t.Errorf("Bad group count: expcected %d, actual %d", 0, len(subscription.GroupIDs))
	}

}
