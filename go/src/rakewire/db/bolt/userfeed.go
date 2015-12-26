package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
)

// UserFeedGetAllByUser retrieves the userfeeds belonging to the user with the Feed populated.
func (z *Service) UserFeedGetAllByUser(userID uint64) ([]*m.UserFeed, error) {

	var result []*m.UserFeed

	// define index keys
	uf := &m.UserFeed{}
	uf.UserID = userID
	minKeys := uf.IndexKeys()[m.UserFeedIndexUser]
	uf.UserID = userID + 1
	nxtKeys := uf.IndexKeys()[m.UserFeedIndexUser]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserFeedEntity)).Bucket([]byte(m.UserFeedIndexUser))
		bUserFeed := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserFeedEntity))
		bFeed := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.FeedEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, bUserFeed); ok {
				uf := &m.UserFeed{}
				if err := uf.Deserialize(data); err != nil {
					return err
				}
				if data, ok := kvGet(uf.FeedID, bFeed); ok {
					f := &m.Feed{}
					if err := f.Deserialize(data); err != nil {
						return err
					}
					uf.Feed = f
					result = append(result, uf)
				}
			}

		}

		return nil

	})

	return result, err

}

// UserFeedSave saves a user to the database.
func (z *Service) UserFeedSave(userfeed *m.UserFeed) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {
		return kvSave(m.UserFeedEntity, userfeed, tx)
	})

	return err

}
