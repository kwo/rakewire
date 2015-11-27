package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"testing"
	"time"
)

func TestSerialization(t *testing.T) {

	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	// start testing
	fl := m.NewFeedLog("0000-FEED-ID")
	fl.ContentLength = 50
	fl.Duration = 6 * time.Millisecond
	fl.IsUpdated = true
	fl.Result = "OK"
	fl.StartTime = time.Now()

	// marshal
	database.Lock()
	err := database.db.Update(func(tx *bolt.Tx) error {
		return Put(fl, tx)
	})
	database.Unlock()
	assertNoError(t, err)

	// // view out of curiosity
	// err = database.db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte(bucketData)).Bucket([]byte(bucketFeedLog)) // works
	// 	c := b.Cursor()
	// 	for k, v := c.First(); k != nil; k, v = c.Next() {
	// 		t.Logf("FeedLog: %s: %s", k, v)
	// 	} // for
	// 	return nil
	// })
	// assertNoError(t, err)

	// unmarshal
	fl2 := &m.FeedLog{
		ID: fl.ID,
	}
	err = database.db.View(func(tx *bolt.Tx) error {
		return Get(fl2, tx)
	})
	assertNoError(t, err)

	// compare
	assertEqual(t, fl.ID, fl2.ID)
	assertEqual(t, fl.FeedID, fl2.FeedID)
	assertEqual(t, fl.ContentLength, fl2.ContentLength)
	assertEqual(t, fl.Duration, fl2.Duration)
	assertEqual(t, fl.IsUpdated, fl2.IsUpdated)
	assertEqual(t, fl.Result, fl2.Result)
	assertEqual(t, fl.StartTime.UTC(), fl2.StartTime)
	assertEqual(t, fl.LastUpdated, fl2.LastUpdated)
	// zero values are not saved
	assertEqual(t, 0, fl2.StatusCode)
	assertEqual(t, false, fl2.UsesGzip)
	assertEqual(t, "", fl2.ETag)
	assertEqual(t, time.Time{}, fl2.LastUpdated)

	// modify and resave
	fl2.IsUpdated = false
	database.Lock()
	err = database.db.Update(func(tx *bolt.Tx) error {
		return Put(fl2, tx)
	})
	database.Unlock()
	assertNoError(t, err)

	// assert key is not present
	err = database.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(bucketFeedLog)) // works
		value := b.Get([]byte(fmt.Sprintf("%s%s%s", fl2.ID, chSep, "IsUpdated")))
		if value != nil {
			t.Error("value must be nil")
		}
		return nil
	})
	assertNoError(t, err)

}

func TestQuery(t *testing.T) {

	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	result := []*m.FeedLog{}
	add := func() interface{} {
		fl := m.NewFeedLog("0000-FEED-ID")
		result = append(result, fl)
		return fl
	}

	err := database.db.View(func(tx *bolt.Tx) error {
		return Query("FeedLog", empty, []interface{}{" "}, []interface{}{" "}, add, tx)
	})
	assertNoError(t, err)
	//assertEqual(t, 5, len(result))

}
