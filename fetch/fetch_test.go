package fetch

import (
	"github.com/stretchr/testify/require"
	"os"
	m "rakewire.com/model"
	"testing"
)

func TestFetch(t *testing.T) {

	//t.SkipNow()

	r, err := os.Open("../test/feedlist.txt")
	require.Nil(t, err)
	require.NotNil(t, r)
	feeds := m.ParseListToFeeds(r)
	r.Close()
	require.NotNil(t, feeds)

	err = Fetch(feeds)
	require.Nil(t, err)

}
