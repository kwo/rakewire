package logging

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPackageNames(t *testing.T) {
	type x struct{}
	assert.Equal(t, "rakewire.com/logging", reflect.TypeOf(x{}).PkgPath())
}
