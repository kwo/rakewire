package bolt

import (
	"bufio"
	//"fmt"
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

	urls, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, urls)

	feeds := db.NewFeeds()
	for _, url := range urls {
		feeds.Add(db.NewFeed(url))
	}
	assert.Equal(t, len(urls), feeds.Size())

	database := Database{}
	err = database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	updateCount, err := database.SaveFeeds(feeds)
	require.Nil(t, err)
	assert.Equal(t, feeds.Size(), updateCount)

	feeds2, err := database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

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
