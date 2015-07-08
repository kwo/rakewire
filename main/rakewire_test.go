package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHomeDirectory(t *testing.T) {
	assert.Equal(t, "/Users/karl", getHomeDirectory())
}

func TestConfigFileLocation(t *testing.T) {
	assert.Equal(t, "/Users/karl/.config/rakewire/config.yaml", getConfigFileLocation())
}

func TestConfig(t *testing.T) {
	assert.NotNil(t, getConfig())
}
