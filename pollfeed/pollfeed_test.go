package pollfeed

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	cfg := &Configuration{}

	qf := NewService(cfg)

	qf.Start()
	require.Equal(t, true, qf.IsRunning())
	qf.Stop()
	assert.Equal(t, false, qf.IsRunning())

}
