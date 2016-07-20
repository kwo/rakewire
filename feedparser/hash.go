package feedparser

import (
	"crypto/md5"
	"encoding/hex"
)

// HashFeed assigns a hash of the entry to entry.ID only if the entry lacks an ID.
func HashFeed(f *Feed) {
	for _, entry := range f.Entries {
		if isEmpty(entry.ID) {
			entry.ID = HashEntry(entry)
		}
	}
}

// HashEntry calculates a unique signature for an entry, used if the entry lacks a unique ID/GUID.
func HashEntry(entry *Entry) string {

	// Needs to be reproducible so that on future parsings of the same entry, the same ID will ge generated.
	// Also needs to have enough variables to be unique within the feed.

	hash := md5.New()
	hash.Write([]byte(entry.Content))
	hash.Write([]byte(entry.Summary))
	hash.Write([]byte(entry.Title))
	hash.Write([]byte(entry.Links[linkAlternate]))
	return hex.EncodeToString(hash.Sum(nil))

}
