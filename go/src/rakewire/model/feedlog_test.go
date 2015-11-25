package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFeedLog(t *testing.T) {

	fl := NewFeedLog("123")
	require.NotNil(t, fl)
	assert.Equal(t, "123", fl.FeedID)
	assert.NotNil(t, fl.ID)
	assert.Equal(t, 36, len(fl.ID))

}
