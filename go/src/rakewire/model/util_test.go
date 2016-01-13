package model

import (
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestParseFeedsFromFile(t *testing.T) {

	t.Parallel()

	feeds, err := ParseFeedsFromFile("../../../test/feedlistmini.txt")
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 10, len(feeds))
	assertEqual(t, "http://www.addrup.de/feed.xml", feeds[0].URL)

}

func TestParseFeedsFromFileError(t *testing.T) {

	t.Parallel()

	_, err := ParseFeedsFromFile("../../../test/feedlistmini2.txt")
	assertNotNil(t, err)

}

func openDatabase() (*bolt.DB, error) {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := bolt.Open(f.Name(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil

}

func closeDatabase(db *bolt.DB) error {

	if err := db.Close(); err != nil {
		return err
	}

	if err := os.Remove(db.Path()); err != nil {
		return err
	}

	return nil

}
