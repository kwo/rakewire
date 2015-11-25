package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFeedsFromFile(t *testing.T) {

	feeds, err := ParseFeedsFromFile("../../../test/feedlistmini.txt")
	require.Nil(t, err)
	require.NotNil(t, feeds)
	assert.Equal(t, 10, len(feeds))
	assert.Equal(t, "http://www.addrup.de/feed.xml", feeds[0].URL)

	logger.Debugf("feeds: %d\n", len(feeds))

}

func TestParseFeedsFromFileError(t *testing.T) {

	_, err := ParseFeedsFromFile("../../../test/feedlistmini2.txt")
	require.NotNil(t, err)

}
