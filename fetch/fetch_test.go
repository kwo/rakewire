package fetch

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetch(t *testing.T) {

	t.SkipNow()

	var err = Fetch("../test/feedlist.txt")
	require.Nil(t, err)

}
