package config

import (
	"fmt"
	"testing"
)

const (
	testConfigFile = "../../../test/config.yaml"
)

func TestConfiguration(t *testing.T) {

	c := Configuration{}
	err := c.LoadFromFile(testConfigFile)
	assertNoError(t, err)

	assertNotNil(t, c.Httpd)

	assertEqual(t, "", c.Httpd.Address)
	assertEqual(t, 4444, c.Httpd.Port)
	assertEqual(t, ":4444", fmt.Sprintf("%s:%d", c.Httpd.Address, c.Httpd.Port))

	assertEqual(t, "/Users/karl/.rakewire/data.db", c.Database.Location)

	assertNotNil(t, c.Fetcher)
	assertEqual(t, 10, c.Fetcher.Workers)

}

func assertNoError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("Error: %s", e.Error())
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
