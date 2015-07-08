package bolt

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	"strings"
	"testing"
)

const (
	feedFile     = "../../test/feedlist.txt"
	databaseFile = "../../test/test.db"
)

func TestFeeds(t *testing.T) {

	feeds, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, feeds)

	var feedinfos []*db.FeedInfo
	for _, url := range feeds {
		feedinfos = append(feedinfos, db.NewFeedInfo(url))
	}
	assert.Equal(t, len(feeds), len(feedinfos))

	db := Database{}
	err = db.init(databaseFile)
	require.Nil(t, err)

	updateCount, err := db.saveFeeds(feedinfos)
	require.Nil(t, err)
	assert.Equal(t, len(feeds), updateCount)

	feedmap, err := db.getFeeds()
	require.Nil(t, err)
	require.NotNil(t, feedmap)

	err = db.destroy()
	assert.Nil(t, err)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

	assert.Equal(t, len(feedmap), len(feeds))
	for k, v := range feedmap {
		fmt.Printf("Feed %s: %v\n", k, v.URL)
	}

}

func readFile(feedfile string) ([]string, error) {

	var result []string

	f, err1 := os.Open(feedfile)
	if err1 != nil {
		return nil, err1
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			result = append(result, url)
		}
	}
	f.Close()

	return result, nil

}
