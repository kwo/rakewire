package modelng

import (
	"bytes"
	"time"
)

// T transmissions all transmission database methods
var T = &transmissionStore{}

type transmissionStore struct{}

func (z *transmissionStore) Delete(id string, tx Transaction) error {
	return delete(entityTransmission, id, tx)
}

func (z *transmissionStore) Get(id string, tx Transaction) *Transmission {
	bData := tx.Bucket(bucketData, entityTransmission)
	if data := bData.Get([]byte(id)); data != nil {
		transmission := &Transmission{}
		if err := transmission.decode(data); err == nil {
			return transmission
		}
	}
	return nil
}

func (z *transmissionStore) GetForFeed(feedID string, since time.Duration, tx Transaction) Transmissions {
	// index Transmission FeedTime = FeedID|StartTime : TransmissionID
	transmissions := Transmissions{}
	now := time.Now().Truncate(time.Second)
	min := []byte(keyEncode(feedID, keyEncodeTime(now.Add(-since))))
	max := []byte(keyEncode(feedID, keyEncodeTime(now)))
	b := tx.Bucket(bucketIndex, entityTransmission, indexTransmissionFeedTime)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		transmissionID := string(v)
		if transmission := T.Get(transmissionID, tx); transmission != nil {
			transmissions = append(transmissions, transmission)
		}
	}
	transmissions.Reverse()
	return transmissions
}

func (z *transmissionStore) GetLast(tx Transaction) *Transmission {
	// index Transmission Time = StartTime|TransmissionID : TransmissionID
	b := tx.Bucket(bucketData, entityTransmission)
	c := b.Cursor()
	if k, _ := c.Last(); k != nil {
		transmissionID := string(k)
		return T.Get(transmissionID, tx)
	}
	return nil
}

func (z *transmissionStore) GetRange(maxTime time.Time, since time.Duration, tx Transaction) Transmissions {
	// index Transmission Time = StartTime|TransmissionID : TransmissionID
	transmissions := Transmissions{}
	minTime := maxTime.Add(-since)
	min := []byte(keyEncodeTime(minTime))
	max := []byte(keyEncodeTime(maxTime))
	b := tx.Bucket(bucketIndex, entityTransmission, indexTransmissionTime)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		transmissionID := string(v)
		if transmission := T.Get(transmissionID, tx); transmission != nil {
			transmissions = append(transmissions, transmission)
		}
	}
	transmissions.Reverse()
	return transmissions
}

func (z *transmissionStore) New(feedID string) *Transmission {
	return &Transmission{
		FeedID: feedID,
	}
}

func (z *transmissionStore) Save(transmission *Transmission, tx Transaction) error {
	return save(entityTransmission, transmission, tx)
}
