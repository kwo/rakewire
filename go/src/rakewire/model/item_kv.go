package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// index names
const (
	itemEntity    = "Item"
	itemIndexGUID = "GUID"
)

const (
	itemID      = "ID"
	itemGUID    = "GUID"
	itemFeedID  = "FeedID"
	itemCreated = "Created"
	itemUpdated = "Updated"
	itemURL     = "URL"
	itemAuthor  = "Author"
	itemTitle   = "Title"
	itemContent = "Content"
)

var (
	itemAllFields = []string{
		itemID, itemGUID, itemFeedID, itemCreated, itemUpdated, itemURL, itemAuthor, itemTitle, itemContent,
	}
	itemAllIndexes = []string{
		itemIndexGUID,
	}
)

// Items is a collection of Item elements
type Items []*Item

func (z Items) Len() int      { return len(z) }
func (z Items) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Items) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Items) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Items) First() *Item { return z[0] }

// Reverse reverses the order of the collection
func (z Items) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Item) getID() uint64 {
	return z.ID
}

// setID sets the primary key of the object.
func (z *Item) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Item) clear() {
	z.ID = 0
	z.GUID = ""
	z.FeedID = 0
	z.Created = time.Time{}
	z.Updated = time.Time{}
	z.URL = ""
	z.Author = ""
	z.Title = ""
	z.Content = ""

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Item) serialize(flags ...bool) Record {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[itemID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.GUID != "" {
		result[itemGUID] = z.GUID
	}

	if flagNoZeroCheck || z.FeedID != 0 {
		result[itemFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if flagNoZeroCheck || !z.Created.IsZero() {
		result[itemCreated] = z.Created.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[itemUpdated] = z.Updated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.URL != "" {
		result[itemURL] = z.URL
	}

	if flagNoZeroCheck || z.Author != "" {
		result[itemAuthor] = z.Author
	}

	if flagNoZeroCheck || z.Title != "" {
		result[itemTitle] = z.Title
	}

	if flagNoZeroCheck || z.Content != "" {
		result[itemContent] = z.Content
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Item) deserialize(values Record, flags ...bool) error {
	flagUnknownCheck := len(flags) > 0 && flags[0]

	var errors []error
	var missing []string
	var unknown []string

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(itemID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, itemID)
	}

	z.GUID = values[itemGUID]

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(itemFeedID, values, errors)

	if !(z.FeedID != 0) {
		missing = append(missing, itemFeedID)
	}

	z.Created = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(itemCreated, values, errors)

	z.Updated = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(itemUpdated, values, errors)

	z.URL = values[itemURL]

	z.Author = values[itemAuthor]

	z.Title = values[itemTitle]

	z.Content = values[itemContent]

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, itemAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(itemEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Item) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[itemIndexGUID] = []string{

		data[itemFeedID],

		data[itemGUID],
	}

	return result
}

func newItemID(tx Transaction) (string, error) {
	return kvNextID(itemEntity, tx)
}

// GroupByGUID groups elements in the Items collection by GUID
func (z Items) GroupByGUID() map[string]*Item {
	result := make(map[string]*Item)
	for _, item := range z {
		result[item.GUID] = item
	}
	return result
}

// GroupAllByFeedID groups collections of elements in Items by FeedID
func (z Items) GroupAllByFeedID() map[uint64]Items {
	result := make(map[uint64]Items)
	for _, item := range z {
		a := result[item.FeedID]
		a = append(a, item)
		result[item.FeedID] = a
	}
	return result
}
