package fetch

import (
	"github.com/kwo/rakewire/model"
	"testing"
)

func TestInterfaceService(t *testing.T) {

	var s model.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}
