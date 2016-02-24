package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"sort"
	"strings"
)

// index names
const (
	groupEntity         = "Group"
	groupIndexUserGroup = "UserGroup"
)

const (
	groupID     = "ID"
	groupUserID = "UserID"
	groupName   = "Name"
)

var (
	groupAllFields = []string{
		groupID, groupUserID, groupName,
	}
	groupAllIndexes = []string{
		groupIndexUserGroup,
	}
)

// Groups is a collection of Group elements
type Groups []*Group

func (z Groups) Len() int      { return len(z) }
func (z Groups) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Groups) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Groups) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Groups) First() *Group { return z[0] }

// Reverse reverses the order of the collection
func (z Groups) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Group) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Group) clear() {
	z.ID = empty
	z.UserID = empty
	z.Name = empty

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Group) serialize(flags ...bool) Record {

	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(Record)

	if flagNoZeroCheck || z.ID != empty {
		result[groupID] = z.ID
	}

	if flagNoZeroCheck || z.UserID != empty {
		result[groupUserID] = z.UserID
	}

	if flagNoZeroCheck || z.Name != empty {
		result[groupName] = z.Name
	}

	return result

}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Group) deserialize(values Record, flags ...bool) error {

	flagUnknownCheck := len(flags) > 0 && flags[0]
	z.clear()

	var errors []error
	var missing []string
	var unknown []string

	z.ID = values[groupID]
	if !(z.ID != empty) {
		missing = append(missing, groupID)
	}

	z.UserID = values[groupUserID]
	if !(z.UserID != empty) {
		missing = append(missing, groupUserID)
	}

	z.Name = values[groupName]
	if !(z.Name != empty) {
		missing = append(missing, groupName)
	}

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, groupAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}

	return newDeserializationError(groupEntity, errors, missing, unknown)

}

// serializeIndexes returns all index records
func (z *Group) serializeIndexes() map[string]Record {

	result := make(map[string]Record)
	data := z.serialize(true)
	var keys []string

	keys = []string{}
	keys = append(keys, data[groupUserID])
	keys = append(keys, strings.ToLower(data[groupName]))
	result[groupIndexUserGroup] = Record{string(kvKeyEncode(keys...)): data[groupID]}

	return result

}

// GroupByID groups elements in the Groups collection by ID
func (z Groups) GroupByID() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.ID] = group
	}
	return result
}

// GroupByName groups elements in the Groups collection by Name
func (z Groups) GroupByName() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.Name] = group
	}
	return result
}
