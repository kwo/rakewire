package model

import (
	"testing"
)

func TestAppVariables(t *testing.T) {
	t.SkipNow()
	if BuildHash == "" {
		t.Error("Build hash not set")
	}
	if BuildTime == "" {
		t.Error("Build time not set")
	}
	if Version == "" {
		t.Error("Version not set")
	}
}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Fatal("Expected nil value")
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("Expected not nil value")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Not equal: expected %v, actual %v", a, b)
	}
}
