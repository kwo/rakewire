package model

import (
	"fmt"
	"testing"
	"time"
)

func TestKVDelete(t *testing.T) {

	t.Parallel()

	d := openTestDatabase(t)
	defer closeTestDatabase(t, d)

	const URL = "http://localhost/"

	err := d.Update(func(tx Transaction) error {
		for i := 0; i < 5; i++ {
			f := NewFeed(fmt.Sprintf("%s%d", URL, i))
			f.ETag = fmt.Sprintf("%s%d", "Etag-", i)
			f.Title = fmt.Sprintf("%s%d", "Title-", i)
			f.Notes = fmt.Sprintf("%s%d", "Notes-", i)
			f.LastModified = time.Now()
			_, err := f.Save(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error saving new feeds: %s", err.Error())
	}

	err = d.Update(func(tx Transaction) error {
		f, err := FeedByURL("http://localhost/2", tx)
		if err != nil {
			return err
		}
		return f.Delete(tx)
	})
	if err != nil {
		t.Fatalf("Error deleting feed: %s", err.Error())
	}

	err = d.Select(func(tx Transaction) error {
		b := tx.Bucket("Data").Bucket("Feed")
		return b.ForEach(func(key, value []byte) error {
			id, err := kvKeyElementID(key, 0)
			if err != nil {
				t.Errorf("Error parsing ID from key: %s", err.Error())
			} else if id == 3 {
				t.Error("Deleted feed ID still present in bucket")
			}
			t.Logf("%s: %s", key, value)
			return nil
		})
	})
	if err != nil {
		t.Errorf("Error viewing feed: %s", err.Error())
	}

}
