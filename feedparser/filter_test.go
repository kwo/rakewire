package feedparser

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestFilter1(t *testing.T) {

	r := strings.NewReader("")
	f := newFilterReader(r)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("Cannot read from filter: %s", err.Error())
	}

	if len(data) != 0 {
		t.Errorf("Expected 0, actual: %d", len(data))
	}

}

func TestFilter2(t *testing.T) {

	r := strings.NewReader("")
	f := newFilterReader(r)

	b := []byte{}
	i, err := f.Read(b)

	if err != nil {
		t.Errorf("Cannot read from filter: %s", err.Error())
	}

	if i != 0 {
		t.Errorf("Expected 0, actual: %d", i)
	}

}

func TestFilter3(t *testing.T) {

	r := strings.NewReader("abcdefg\nhijklmnop\tqrstuv\rwxyz")
	f := newFilterReader(r)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("Cannot read from filter: %s", err.Error())
	}

	if len(data) != 26 {
		t.Errorf("Expected 26, actual: %d", len(data))
	}

	if string(data) != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("Expected alphabet, actual: %s", string(data))
	}

}

func TestFilter4(t *testing.T) {

	b := "1234567890"

	s := b[:5]

	expected := "12345"
	if s != expected {
		t.Errorf("Expected %s, actual %s", expected, s)
	}

}
