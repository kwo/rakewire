package bolt

import (
	"github.com/stretchr/testify/assert"
	"rakewire/db"
	"rakewire/logging"
	"testing"
)

const (
	feedFile     = "../../../../test/feedlistmini.txt"
	databaseFile = "../../../../test/bolt.db"
)

func TestMain(m *testing.M) {

	// initialize logging
	logging.Init(&logging.Configuration{
		Level: "debug",
	})

	logger.Debug("Logging configured")

	m.Run()

}

func TestInterface(t *testing.T) {

	var d db.Database = &Database{}
	assert.NotNil(t, d)

}
