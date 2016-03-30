package model

import (
	"fmt"
	"github.com/murphysean/cache"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	logName = "[db]"
	logInfo = "[INFO]"
)

type lookupFunc func(id ...string) bool

func (z *boltInstance) Check(filename string) error {

	// backup database (rename to backup name)
	var backupFilename string
	if name, err := z.backupDatabase(filename); err == nil {
		backupFilename = name
	} else {
		return err
	}

	// open both old and new databases
	// check schema of both databases (in Open)
	var oldDb Database
	if db, err := z.Open(backupFilename); err == nil {
		oldDb = db
	} else {
		return err
	}
	defer z.Close(oldDb)

	var newDb Database
	if db, err := z.Open(filename); err == nil {
		newDb = db
	} else {
		return err
	}
	defer z.Close(newDb)

	// copy buckets
	if err := z.copyBuckets(oldDb, newDb); err != nil {
		return err
	}

	// TODO: check integrity of each bucket
	log.Printf("%-7s %-7s validating buckets...", logInfo, logName)
	if err := z.checkSubscriptions(newDb); err != nil {
		return err
	}
	log.Printf("%-7s %-7s validating buckets done", logInfo, logName)

	// rebuild indexes
	if err := z.rebuildIndexes(oldDb, newDb); err != nil {
		return err
	}

	return nil

}

func (z *boltInstance) backupDatabase(location string) (string, error) {

	now := time.Now().Truncate(time.Second)
	timestamp := now.Format("20060102150405")

	dir := filepath.Dir(location)
	ext := filepath.Ext(location)
	filename := strings.TrimSuffix(filepath.Base(location), ext)

	newFilename := fmt.Sprintf("%s%s%s-%s%s", dir, string(os.PathSeparator), filename, timestamp, ext)
	err := os.Rename(location, newFilename)

	return newFilename, err

}

func (z *boltInstance) makeFeedCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if f := F.Get(tx, key); f != nil {
			return f, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		feedID := id[0]
		if f, err := c.Get(feedID); f != nil && err == nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) makeGroupCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if g := G.Get(tx, key); g != nil {
			return g, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		userID := id[0]
		groupID := id[1]
		if g, err := c.Get(groupID); g != nil && err == nil {
			return g.(*Group).UserID == userID
		}
		return false
	}

}

func (z *boltInstance) makeUserCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if u := U.Get(tx, key); u != nil {
			return u, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		userID := id[0]
		if u, err := c.Get(userID); u != nil && err == nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) checkSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   %s ...", logInfo, logName, entitySubscription)

	return db.Update(func(tx Transaction) error {

		userExists := z.makeUserCache(tx, 50)
		feedExists := z.makeFeedCache(tx, 500)
		groupExists := z.makeGroupCache(tx, 500)

		badIDs := []string{}
		cleanedSubscriptions := Subscriptions{}

		subscriptions := tx.Bucket(bucketData, entitySubscription)
		c := subscriptions.Cursor()

		subscription := &Subscription{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			subscription.clear()
			if err := subscription.decode(v); err == nil {
				if !userExists(subscription.UserID) {
					log.Printf("    subscription without user: %s (%s %s)", subscription.UserID, subscription.GetID(), subscription.Title)
					badIDs = append(badIDs, subscription.GetID())
				} else if !feedExists(subscription.FeedID) {
					log.Printf("    subscription without feed: %s (%s %s)", subscription.FeedID, subscription.GetID(), subscription.Title)
					badIDs = append(badIDs, subscription.GetID())
				} else {

					// remove invalid groups
					invalidGroupIDs := []string{}
					for _, groupID := range subscription.GroupIDs {
						if !groupExists(subscription.UserID, groupID) {
							log.Printf("    subscription with invalid group: %s (%s %s)", groupID, subscription.GetID(), subscription.Title)
							invalidGroupIDs = append(invalidGroupIDs, groupID)
						}
					}
					if len(invalidGroupIDs) > 0 {
						for _, groupID := range invalidGroupIDs {
							subscription.RemoveGroup(groupID)
						}
						cleanedSubscriptions = append(cleanedSubscriptions, subscription)
					}

					if len(subscription.GroupIDs) == 0 {
						log.Printf("    subscription without groups: %s %s", subscription.GetID(), subscription.Title)
					}

				}
			}
		}

		// remove bad subscriptions
		for _, id := range badIDs {
			if err := S.Delete(tx, id); err != nil {
				return err
			}
		}

		// resave cleaned subscriptions
		for _, subscription := range cleanedSubscriptions {
			if err := S.Save(tx, subscription); err != nil {
				return err
			}
		}

		return nil

	})

}

func (z *boltInstance) copyBuckets(srcDb, dstDb Database) error {

	log.Printf("%-7s %-7s copying buckets...", logInfo, logName)
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				log.Printf("%-7s %-7s   %s ...", logInfo, logName, entityName)
				srcBucket := srcTx.Bucket(bucketData, entityName)
				dstBucket := dstTx.Bucket(bucketData, entityName)
				entity := getObject(entityName)
				c := srcBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if errDecode := entity.decode(v); errDecode == nil {
						if data, errEncode := entity.encode(); errEncode == nil {
							if err := dstBucket.Put(k, data); err != nil {
								return err
							}
						} else {
							log.Printf("%-7s %-7s     error encoding entity (%s): %s", logInfo, logName, k, errEncode.Error())
						}
					} else {
						log.Printf("%-7s %-7s     error decoding entity (%s): %s", logInfo, logName, k, errDecode.Error())
					}
				} // cursor
			} // entities
			return nil
		}) // update
	}) // select

	if err != nil {
		return err
	}

	log.Printf("%-7s %-7s copying buckets done", logInfo, logName)

	return nil

}

func (z *boltInstance) rebuildIndexes(srcDb, dstDb Database) error {

	log.Printf("%-7s %-7s rebuild indexes...", logInfo, logName)
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				log.Printf("%-7s %-7s   %s ...", logInfo, logName, entityName)
				srcBucket := srcTx.Bucket(bucketData, entityName)
				dstIndex := dstTx.Bucket(bucketIndex, entityName)
				entity := getObject(entityName)
				c := srcBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if err := entity.decode(v); err == nil {
						// save new indexes
						for indexName, indexKeys := range entity.indexes() {
							bIndex := dstIndex.Bucket(indexName)
							key := keyEncode(indexKeys...)
							value := entity.GetID()
							if err := bIndex.Put([]byte(key), []byte(value)); err != nil {
								return err
							}
						} // indexes
					} else {
						// This error should never occur since the entity was freshly encoded when copying the data.
						log.Printf("%-7s %-7s     error decoding entity (%s): %s", logInfo, logName, k, err.Error())
						// Note: just reporting the error means that there is an entity in data without an index
					}
				} // cursor
			} // entities
			return nil
		}) // update
	}) // select

	if err != nil {
		return err
	}

	log.Printf("%-7s %-7s rebuild indexes done", logInfo, logName)

	return nil

}
