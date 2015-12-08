package httpd

import (
	"io/ioutil"
	"testing"
)

func TestStatic(t *testing.T) {

	fs := Dir(true, "/public")
	if fs == nil {
		t.Fatal("fs is nil")
	}

	f, err := fs.Open("/robots.txt")
	if err != nil {
		t.Fatalf("Error opening file: %s", err.Error())
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error opening file: %s", err.Error())
	}

	t.Logf("%v, %v\n", err, string(data))

}

func TestStaticBogus(t *testing.T) {

	fs := Dir(true, "/public")
	if fs == nil {
		t.Fatal("fs is nil")
	}

	f, err := fs.Open("/robots2.txt")
	if err == nil {
		t.Error("Expected error, none returned")
	} else if err.Error() != "file does not exist" {
		t.Errorf("Error message incorrect, actual: %s", err.Error())
	}
	if f != nil {
		t.Errorf("Expected nil file, actual: %v", f)
	}

}
