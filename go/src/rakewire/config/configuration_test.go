package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testConfigFile = "../../../test/config.yaml"
)

func TestConfiguration(t *testing.T) {

	c := Configuration{}
	err := c.LoadFromFile(testConfigFile)
	require.Nil(t, err)

	require.NotNil(t, c.Httpd)

	assert.Equal(t, "", c.Httpd.Address)
	assert.Equal(t, 4444, c.Httpd.Port)
	assert.Equal(t, ":4444", fmt.Sprintf("%s:%d", c.Httpd.Address, c.Httpd.Port))

	assert.Equal(t, "/Users/karl/.rakewire/data.db", c.Database.Location)

	require.NotNil(t, c.Fetcher)
	assert.Equal(t, 10, c.Fetcher.Workers)

}
