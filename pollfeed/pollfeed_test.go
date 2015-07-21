package pollfeed

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	cfg := &Configuration{}

	pf := NewService(cfg, nil)

	pf.Start()
	require.Equal(t, true, pf.IsRunning())
	pf.Stop()
	assert.Equal(t, false, pf.IsRunning())

}
