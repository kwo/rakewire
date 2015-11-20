package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire/db"
	m "rakewire/model"
	"testing"
	"time"
)

func TestSerialization(t *testing.T) {

	// init db
	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// start testing
	fl := m.NewFeedLog("0000-FEED-ID")
	fl.ContentLength = 50
	fl.Duration = 6 * time.Millisecond
	fl.IsUpdated = true
	fl.Result = "OK"
	fl.StartTime = time.Now().UTC().Truncate(time.Millisecond)

	// marshal
	database.Lock()
	err = database.db.Update(func(tx *bolt.Tx) error {
		return marshal(fl, tx)
	})
	database.Unlock()
	assert.Nil(t, err)

	// view out of curiosity
	err = database.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeedLog)) // works
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			logger.Debugf("FeedLog: %s: %s", k, v)
		} // for
		return nil
	})
	assert.Nil(t, err)

	// unmarshal
	fl2 := &m.FeedLog{
		ID: fl.ID,
	}
	err = database.db.View(func(tx *bolt.Tx) error {
		return unmarshal(fl2, tx)
	})
	assert.Nil(t, err)

	// compare
	assert.Equal(t, fl.ID, fl2.ID)
	assert.Equal(t, fl.FeedID, fl2.FeedID)
	assert.Equal(t, fl.ContentLength, fl2.ContentLength)
	assert.Equal(t, fl.Duration, fl2.Duration)
	assert.Equal(t, fl.IsUpdated, fl2.IsUpdated)
	assert.Equal(t, fl.Result, fl2.Result)
	assert.Equal(t, fl.StartTime, fl2.StartTime)
	assert.Equal(t, fl.Updated, fl2.Updated)
	// zero values are not saved
	assert.Equal(t, 0, fl2.StatusCode)
	assert.Equal(t, false, fl2.UsesGzip)
	assert.Equal(t, "", fl2.ETag)
	assert.Equal(t, time.Time{}, fl2.Updated)

	// cleanup
	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}

func TestMetadataPrimaryKey(t *testing.T) {

	type metatest struct {
		ID  string
		Key int `db:"primary-key"`
	}

	mt := &metatest{
		ID:  "1",
		Key: 2,
	}

	meta, err := getMetadata(mt)
	require.Nil(t, err)
	require.NotNil(t, meta)

	assert.Equal(t, "metatest", meta.name)
	assert.Equal(t, "2", meta.key)
	assert.NotNil(t, meta.value)
	assert.Equal(t, 0, len(meta.index))

}

func TestMetadataPrimaryKeyDefault(t *testing.T) {

	type metatest struct {
		ID  string
		Key int
	}

	mt := &metatest{
		ID:  "1",
		Key: 2,
	}

	meta, err := getMetadata(mt)
	require.Nil(t, err)
	require.NotNil(t, meta)

	assert.Equal(t, "metatest", meta.name)
	assert.Equal(t, "1", meta.key)
	assert.NotNil(t, meta.value)
	assert.Equal(t, 0, len(meta.index))

}

func TestMetadataPrimaryKeyEmpty(t *testing.T) {

	type metatest struct {
		ID  string
		Key int
	}

	mt := &metatest{}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Empty primary key for metatest.", err.Error())

}

func TestMetadataPrimaryKeyEmptyInteger(t *testing.T) {

	type metatest struct {
		ID  string
		Key int `db:"primary-key"`
	}

	mt := &metatest{}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Empty primary key for metatest.", err.Error())

}

func TestMetadataPrimaryKeyDuplicate(t *testing.T) {

	type metatest struct {
		ID  string `db:"primary-key"`
		Key int    `db:"primary-key"`
	}

	mt := &metatest{
		ID:  "1",
		Key: 2,
	}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Duplicate primary key defined for metatest.", err.Error())

}

func TestMetadataIndexes(t *testing.T) {

	type metatest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:1"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &metatest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	meta, err := getMetadata(mt)
	require.Nil(t, err)
	require.NotNil(t, meta)

	assert.Equal(t, "metatest", meta.name)
	assert.Equal(t, "1", meta.key)
	assert.NotNil(t, meta.value)
	assert.Equal(t, 2, len(meta.index))
	assert.Nil(t, meta.index["bogusname"])
	assert.NotNil(t, meta.index["Name"])
	assert.NotNil(t, meta.index["URLTitle"])
	assert.Equal(t, 1, len(meta.index["Name"]))
	assert.Equal(t, "name", meta.index["Name"][0])
	assert.Equal(t, 2, len(meta.index["URLTitle"]))
	assert.Equal(t, "url", meta.index["URLTitle"][0])
	assert.Equal(t, "title", meta.index["URLTitle"][1])

}

func TestMetadataIndexesInvalidPosition(t *testing.T) {

	type metatest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:a"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &metatest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Index position is not an integer: indexName:a.", err.Error())

}

func TestMetadataIndexesInvalidDefinition(t *testing.T) {

	type metatest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:1"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle1"`
	}

	mt := &metatest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Invalid index definition: indexURLTitle1.", err.Error())

}

func TestMetadataIndexesZeroPosition(t *testing.T) {

	type metatest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:0"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &metatest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	meta, err := getMetadata(mt)
	require.NotNil(t, err)
	require.Nil(t, meta)

	assert.Equal(t, "Index positions are one-based: indexName:0.", err.Error())

}
