package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"sort"
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
func (z *Item) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Item) clear() {
	z.ID = ""
	z.GUID = ""
	z.FeedID = ""
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

	if flagNoZeroCheck || z.ID != "" {
		result[itemID] = z.ID
	}

	if flagNoZeroCheck || z.GUID != "" {
		result[itemGUID] = z.GUID
	}

	if flagNoZeroCheck || z.FeedID != "" {
		result[itemFeedID] = z.FeedID
	}

	if flagNoZeroCheck || !z.Created.IsZero() {
		result[itemCreated] = z.Created.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[itemUpdated] = z.Updated.UTC().Format(fmtTime)
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

	z.ID = values[itemID]

	if !(z.ID != "") {
		missing = append(missing, itemID)
	}

	z.GUID = values[itemGUID]

	z.FeedID = values[itemFeedID]

	if !(z.FeedID != "") {
		missing = append(missing, itemFeedID)
	}

	z.Created = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(fmtTime, value)
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
			t, err := time.Parse(fmtTime, value)
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

// serializeIndexes returns all index records
func (z *Item) serializeIndexes() map[string]Record {

	result := make(map[string]Record)

	data := z.serialize(true)

	var keys []string

	keys = []string{}

	keys = append(keys, data[itemFeedID])

	keys = append(keys, data[itemGUID])

	result[itemIndexGUID] = Record{string(kvKeyEncode(keys...)): data[itemID]}

	return result
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
func (z Items) GroupAllByFeedID() map[string]Items {
	result := make(map[string]Items)
	for _, item := range z {
		a := result[item.FeedID]
		a = append(a, item)
		result[item.FeedID] = a
	}
	return result
}
