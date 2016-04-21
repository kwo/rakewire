package model

import (
	"encoding/json"
)

const (
	entityGroup        = "Group"
	indexGroupUserName = "UserName"
)

var (
	indexesGroup = []string{
		indexGroupUserName,
	}
)

// Group defines an item status for a user
type Group struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

// GetID returns the unique ID for the object
func (z *Group) GetID() string {
	return z.ID
}

func (z *Group) clear() {
	z.ID = empty
	z.UserID = empty
	z.Name = empty
}

func (z *Group) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Group) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Group) hasIncrementingID() bool {
	return true
}

func (z *Group) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexGroupUserName] = []string{z.UserID, z.Name}
	return result
}

func (z *Group) setID(tx Transaction) error {
	id, err := tx.Bucket(bucketData, entityGroup).NextID()
	if err != nil {
		return err
	}
	z.ID = keyEncodeUint(id)
	return nil
}

// Groups is a collection of Group elements
type Groups []*Group

// ByID groups elements in the Groups collection by ID
func (z Groups) ByID() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.ID] = group
	}
	return result
}

// ByName groups elements in the Groups collection by Name
func (z Groups) ByName() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.Name] = group
	}
	return result
}

func (z *Groups) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Groups) encode() ([]byte, error) {
	return json.Marshal(z)
}
