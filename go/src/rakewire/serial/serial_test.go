package serial

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/logging"
	"testing"
)

func TestMain(m *testing.M) {

	// initialize logging
	logging.Init(&logging.Configuration{
		Level: "debug",
	})

	logger.Debug("Logging configured")

	m.Run()

}

func TestDecode(t *testing.T) {

	type object struct {
		ID string
	}

	o := &object{}
	data := map[string]string{
		"ID": "hello",
	}

	err := Decode(o, data)
	require.Nil(t, err)

	assert.Equal(t, data["ID"], o.ID)

}
