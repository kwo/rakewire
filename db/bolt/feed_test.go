package bolt

import (
	"bufio"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	m "rakewire.com/model"
	"strings"
	"testing"
)

const (
	feedFile     = "../../test/feedlist.txt"
	databaseFile = "../../test/test.db"
)

func TestFeeds(t *testing.T) {

	urls, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, urls)

	feeds := m.NewFeeds()
	for _, url := range urls {
		feeds.Add(m.NewFeed(url))
	}
	assert.Equal(t, len(urls), feeds.Size())

	db := Database{}
	err = db.Open(&m.DatabaseConfiguration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	updateCount, err := db.SaveFeeds(feeds)
	require.Nil(t, err)
	assert.Equal(t, feeds.Size(), updateCount)

	feeds2, err := db.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)

	err = db.Close()
	assert.Nil(t, err)
	assert.Nil(t, db.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

	assert.Equal(t, feeds2.Size(), feeds.Size())
	// for k, v := range feedmap {
	// 	fmt.Printf("Feed %s: %v\n", k, v.URL)
	// }

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
