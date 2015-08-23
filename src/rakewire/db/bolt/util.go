package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

func formatFeedLogKey(id string, dt *time.Time) string {
	if dt == nil {
		return fmt.Sprintf("%s!", id)
	}
	return fmt.Sprintf("%s!%s", id, formatTimeKey(*dt))
}

func fetchKey(f *m.Feed) string {
	return fmt.Sprintf("%s!%s", formatTimeKey(f.NextFetch), f.ID)
}

func formatTimeKey(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func formatMaxTime(t time.Time) string {
	return formatTimeKey(t) + "#"
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
