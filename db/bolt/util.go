package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	"rakewire.com/db"
	"time"
)

func fetchKey(f *db.Feed) string {
	return fmt.Sprintf("%s!%s", formatFetchTime(*f.GetNextFetchTime()), f.ID)
}

func formatFetchTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func formatMaxTime(t time.Time) string {
	return formatFetchTime(t) + "#"
}

func (z *Database) checkIndexForEntries(indexName string, value string, threshold int) error {

	var result []string
	z.db.View(func(tx *bolt.Tx) error {
		i := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(indexName))
		result = z.findAllKeysForValue(i, value)
		return nil
	})

	if len(result) > threshold {
		logger.Printf("multiple keys for %s: %s\n", value, result)
		return bolt.ErrInvalid
	}

	return nil

}

func (z *Database) findAllKeysForValue(b *bolt.Bucket, value string) []string {
	var result []string
	b.ForEach(func(k []byte, v []byte) error {
		if string(v) == value {
			result = append(result, string(k))
		}
		return nil
	})
	return result
}
