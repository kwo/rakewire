package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pborman/uuid"
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
	fl := &m.FeedLog{
		ID:            uuid.NewUUID().String(),
		FeedID:        "0000-FEED-ID",
		ContentLength: 50,
		Duration:      6 * time.Millisecond,
		IsUpdated:     true,
		Result:        "OK",
		StartTime:     time.Now().UTC().Truncate(time.Millisecond),
	}

	// marshal
	database.Lock()
	err = database.db.Update(func(tx *bolt.Tx) error {
		return marshal(fl, tx)
	})
	database.Unlock()
	assert.Nil(t, err)

	// view out of curiosity
	err = database.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(bucketFeedLog)) // works
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

	// modify and resave
	fl2.IsUpdated = false
	database.Lock()
	err = database.db.Update(func(tx *bolt.Tx) error {
		return marshal(fl2, tx)
	})
	database.Unlock()
	assert.Nil(t, err)

	// assert key is not present
	err = database.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(bucketFeedLog)) // works
		value := b.Get([]byte(fmt.Sprintf("%s%s%s", fl2.ID, chSep, "IsUpdated")))
		assert.Nil(t, value)
		return nil
	})
	assert.Nil(t, err)

	// cleanup
	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}

func TestSummaryPrimaryKey(t *testing.T) {

	type summarytest struct {
		ID  string
		Key int `db:"primary-key"`
	}

	mt := &summarytest{
		ID:  "1",
		Key: 2,
	}

	summary, err := getSummary(mt)
	require.Nil(t, err)
	require.NotNil(t, summary)

	assert.Equal(t, "summarytest", summary.name)
	assert.Equal(t, "2", summary.key)
	assert.NotNil(t, summary.value)
	assert.Equal(t, 0, len(summary.index))

}

func TestSummaryPrimaryKeyDefault(t *testing.T) {

	type summarytest struct {
		ID  string
		Key int
	}

	mt := &summarytest{
		ID:  "1",
		Key: 2,
	}

	summary, err := getSummary(mt)
	require.Nil(t, err)
	require.NotNil(t, summary)

	assert.Equal(t, "summarytest", summary.name)
	assert.Equal(t, "1", summary.key)
	assert.NotNil(t, summary.value)
	assert.Equal(t, 0, len(summary.index))

}

func TestSummaryPrimaryKeyEmpty(t *testing.T) {

	type summarytest struct {
		ID  string
		Key int
	}

	mt := &summarytest{}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Empty primary key for summarytest.", err.Error())

}

func TestSummaryPrimaryKeyEmptyInteger(t *testing.T) {

	type summarytest struct {
		ID  string
		Key int `db:"primary-key"`
	}

	mt := &summarytest{}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Empty primary key for summarytest.", err.Error())

}

func TestSummaryPrimaryKeyDuplicate(t *testing.T) {

	type summarytest struct {
		ID  string `db:"primary-key"`
		Key int    `db:"primary-key"`
	}

	mt := &summarytest{
		ID:  "1",
		Key: 2,
	}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Duplicate primary key defined for summarytest.", err.Error())

}

func TestSummaryIndexes(t *testing.T) {

	type summarytest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:1"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &summarytest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	summary, err := getSummary(mt)
	require.Nil(t, err)
	require.NotNil(t, summary)

	assert.Equal(t, "summarytest", summary.name)
	assert.Equal(t, "1", summary.key)
	assert.NotNil(t, summary.value)
	assert.Equal(t, 2, len(summary.index))
	assert.Nil(t, summary.index["bogusname"])
	assert.NotNil(t, summary.index["Name"])
	assert.NotNil(t, summary.index["URLTitle"])
	assert.Equal(t, 1, len(summary.index["Name"]))
	assert.Equal(t, "name", summary.index["Name"][0])
	assert.Equal(t, 2, len(summary.index["URLTitle"]))
	assert.Equal(t, "url", summary.index["URLTitle"][0])
	assert.Equal(t, "title", summary.index["URLTitle"][1])

}

func TestSummaryIndexesInvalidPosition(t *testing.T) {

	type summarytest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:a"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &summarytest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Index position is not an integer: indexName:a.", err.Error())

}

func TestSummaryIndexesInvalidDefinition(t *testing.T) {

	type summarytest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:1"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle1"`
	}

	mt := &summarytest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Invalid index definition: indexURLTitle1.", err.Error())

}

func TestSummaryIndexesZeroPosition(t *testing.T) {

	type summarytest struct {
		ID    string `db:"primary-key"`
		Key   int
		Name  string `db:"indexName:0"`
		Title string `db:"indexURLTitle:2"`
		URL   string `db:"indexURLTitle:1"`
	}

	mt := &summarytest{
		ID:    "1",
		Key:   2,
		Name:  "name",
		Title: "title",
		URL:   "url",
	}

	summary, err := getSummary(mt)
	require.NotNil(t, err)
	require.Nil(t, summary)

	assert.Equal(t, "Index positions are one-based: indexName:0.", err.Error())

}
