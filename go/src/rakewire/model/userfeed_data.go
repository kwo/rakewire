package model

import (
	"bytes"
	"strconv"
)

// UserFeedsByUser retrieves the userfeeds belonging to the user with the Feed populated.
func UserFeedsByUser(userID uint64, tx Transaction) ([]*UserFeed, error) {

	var result []*UserFeed

	// define index keys
	uf := &UserFeed{}
	uf.UserID = userID
	minKeys := uf.IndexKeys()[UserFeedIndexUser]
	uf.UserID = userID + 1
	nxtKeys := uf.IndexKeys()[UserFeedIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserFeedEntity).Bucket(UserFeedIndexUser)
	bUserFeed := tx.Bucket(bucketData).Bucket(UserFeedEntity)
	bFeed := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bUserFeed); ok {
			uf := &UserFeed{}
			if err := uf.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.Deserialize(data); err != nil {
					return nil, err
				}
				uf.Feed = f
				result = append(result, uf)
			}
		}

	}

	return result, nil

}

// UserFeedsByFeed retrieves the userfeeds associated with the feed.
func UserFeedsByFeed(feedID uint64, tx Transaction) ([]*UserFeed, error) {

	var result []*UserFeed

	// define index keys
	uf := &UserFeed{}
	uf.FeedID = feedID
	minKeys := uf.IndexKeys()[UserFeedIndexFeed]
	uf.FeedID = feedID + 1
	nxtKeys := uf.IndexKeys()[UserFeedIndexFeed]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserFeedEntity).Bucket(UserFeedIndexFeed)
	bUserFeed := tx.Bucket(bucketData).Bucket(UserFeedEntity)
	bFeed := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bUserFeed); ok {
			uf := &UserFeed{}
			if err := uf.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(uf.FeedID, bFeed); ok {
				f := &Feed{}
				if err := f.Deserialize(data); err != nil {
					return nil, err
				}
				uf.Feed = f
				result = append(result, uf)
			}
		}

	}

	return result, nil

}

// Delete removes a userfeed from the database.
func (userfeed *UserFeed) Delete(tx Transaction) error {
	return kvDelete(UserFeedEntity, userfeed, tx)
}

// Save saves a user to the database.
func (userfeed *UserFeed) Save(tx Transaction) error {
	return kvSave(UserFeedEntity, userfeed, tx)
}
