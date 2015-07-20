package pump

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	cfg := &Configuration{}

	pump := NewService(cfg)
	chErrors := make(chan error)

	pump.Start(chErrors)
	require.Equal(t, true, pump.IsRunning())
	pump.Stop()
	assert.Equal(t, false, pump.IsRunning())

}
