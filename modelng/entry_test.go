package modelng

import (
	"testing"
)

func TestEntrySetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityEntry); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityEntry]; obj == nil {
		t.Error("missing allEntities entry")
	}

}
