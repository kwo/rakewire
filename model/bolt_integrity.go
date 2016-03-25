package model

import (
	"fmt"
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
				dstBucket := dstTx.Bucket(bucketData, entityName)
				entity := getObject(entityName)
				c := srcBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if err := entity.decode(v); err == nil {
						// save new indexes
						for indexName, indexKeys := range entity.indexes() {
							bIndex := dstBucket.Bucket(bucketIndex, entityName, indexName)
							if err := bIndex.Put([]byte(keyEncode(indexKeys...)), []byte(entity.GetID())); err != nil {
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
