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
		for i := 1; i < 6; i++ {
			f := NewFeed(fmt.Sprintf("%s%d", URL, i))
			f.ETag = fmt.Sprintf("%s%d", "Etag-", i)
			f.Title = fmt.Sprintf("%s%d", "Title-", i)
			f.Notes = fmt.Sprintf("%s%d", "Notes-", i)
			f.LastModified = time.Now()
			//t.Logf("%d: %d %s", i, f.ID, f.URL)
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
		f, err := FeedByURL("http://localhost/3", tx)
		if err != nil {
			return err
		}
		return f.Delete(tx)
	})
	if err != nil {
		t.Fatalf("Error deleting feed: %s", err.Error())
	}

	err = d.Select(func(tx Transaction) error {
		c := tx.Bucket("Data").Bucket("Feed").Cursor()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			id := kvKeyDecode(key)[0]
			if id == kvKeyUintEncode(3) {
				t.Error("Deleted feed ID still present in bucket")
			}
			t.Logf("%s: %s", key, value)
			return nil
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error viewing feed: %s", err.Error())
	}

}

func TestDeserialize(t *testing.T) {

	g1 := NewGroup(kvKeyUintEncode(3), "three")
	g1.ID = kvKeyUintEncode(3)
	values := g1.serialize()

	g2 := &Group{}
	if err := g2.deserialize(values, true); err != nil {
		t.Errorf("deserialization error: %s", err.Error())
	}

	values["uuid"] = "unknown-field"

	g2 = &Group{}
	if err := g2.deserialize(values, true); err == nil {
		t.Error("expected deserialization error, none returned")
	} else if derr, ok := err.(*DeserializationError); ok {

		if len(derr.UnknownFieldnames) != 1 || !isStringInArray("uuid", derr.UnknownFieldnames) {
			t.Error("Expected field uuid in UnknownFieldnames")
		}

	} else {
		t.Error("Invalid error returned, not a Deserialization error")
	}

}
