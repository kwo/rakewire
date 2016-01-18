package fetch

import (
	m "rakewire/model"
	"testing"
)

func TestInterfaceService(t *testing.T) {

	var s m.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}
