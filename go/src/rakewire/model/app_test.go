package model

import (
	"testing"
)

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
