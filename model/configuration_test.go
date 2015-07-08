package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testConfigFile = "../test/config.yaml"
)

func TestConfiguration(t *testing.T) {

	c := Configuration{}
	err := c.LoadFromFile(testConfigFile)
	require.Nil(t, err)

	assert.NotNil(t, c.Httpd)

	assert.Equal(t, 4444, c.Httpd.Port)
	assert.Equal(t, "/Users/karl/static", c.Httpd.WebAppDir)

}
