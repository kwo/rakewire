package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"sort"
	"strings"
	"time"
)

// index names
const (
	feedEntity         = "Feed"
	feedIndexNextFetch = "NextFetch"
	feedIndexURL       = "URL"
)

const (
	feedID            = "ID"
	feedURL           = "URL"
	feedSiteURL       = "SiteURL"
	feedETag          = "ETag"
	feedLastModified  = "LastModified"
	feedLastUpdated   = "LastUpdated"
	feedNextFetch     = "NextFetch"
	feedNotes         = "Notes"
	feedTitle         = "Title"
	feedStatus        = "Status"
	feedStatusMessage = "StatusMessage"
	feedStatusSince   = "StatusSince"
)

var (
	feedAllFields = []string{
		feedID, feedURL, feedSiteURL, feedETag, feedLastModified, feedLastUpdated, feedNextFetch, feedNotes, feedTitle, feedStatus, feedStatusMessage, feedStatusSince,
	}
	feedAllIndexes = []string{
		feedIndexNextFetch, feedIndexURL,
	}
)

// Feeds is a collection of Feed elements
type Feeds []*Feed

func (z Feeds) Len() int      { return len(z) }
func (z Feeds) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Feeds) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Feeds) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Feeds) First() *Feed { return z[0] }

// Reverse reverses the order of the collection
func (z Feeds) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Feed) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Feed) clear() {
	z.ID = ""
	z.URL = ""
	z.SiteURL = ""
	z.ETag = ""
	z.LastModified = time.Time{}
	z.LastUpdated = time.Time{}
	z.NextFetch = time.Time{}
	z.Notes = ""
	z.Title = ""
	z.Status = ""
	z.StatusMessage = ""
	z.StatusSince = time.Time{}

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Feed) serialize(flags ...bool) Record {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != "" {
		result[feedID] = z.ID
	}

	if flagNoZeroCheck || z.URL != "" {
		result[feedURL] = z.URL
	}

	if flagNoZeroCheck || z.SiteURL != "" {
		result[feedSiteURL] = z.SiteURL
	}

	if flagNoZeroCheck || z.ETag != "" {
		result[feedETag] = z.ETag
	}

	if flagNoZeroCheck || !z.LastModified.IsZero() {
		result[feedLastModified] = z.LastModified.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || !z.LastUpdated.IsZero() {
		result[feedLastUpdated] = z.LastUpdated.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || !z.NextFetch.IsZero() {
		result[feedNextFetch] = z.NextFetch.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || z.Notes != "" {
		result[feedNotes] = z.Notes
	}

	if flagNoZeroCheck || z.Title != "" {
		result[feedTitle] = z.Title
	}

	if flagNoZeroCheck || z.Status != "" {
		result[feedStatus] = z.Status
	}

	if flagNoZeroCheck || z.StatusMessage != "" {
		result[feedStatusMessage] = z.StatusMessage
	}

	if flagNoZeroCheck || !z.StatusSince.IsZero() {
		result[feedStatusSince] = z.StatusSince.UTC().Format(fmtTime)
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Feed) deserialize(values Record, flags ...bool) error {
	flagUnknownCheck := len(flags) > 0 && flags[0]

	var errors []error
	var missing []string
	var unknown []string

	z.ID = values[feedID]

	if !(z.ID != "") {
		missing = append(missing, feedID)
	}

	z.URL = values[feedURL]

	if !(z.URL != "") {
		missing = append(missing, feedURL)
	}

	z.SiteURL = values[feedSiteURL]

	z.ETag = values[feedETag]

	z.LastModified = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(feedLastModified, values, errors)

	z.LastUpdated = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(feedLastUpdated, values, errors)

	z.NextFetch = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(feedNextFetch, values, errors)

	z.Notes = values[feedNotes]

	z.Title = values[feedTitle]

	z.Status = values[feedStatus]

	z.StatusMessage = values[feedStatusMessage]

	z.StatusSince = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(feedStatusSince, values, errors)

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, feedAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(feedEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Feed) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[feedIndexNextFetch] = []string{

		data[feedNextFetch],

		data[feedID],
	}

	result[feedIndexURL] = []string{

		strings.ToLower(data[feedURL]),
	}

	return result
}

// serializeIndexes returns all index records
func (z *Feed) serializeIndexes() map[string]Record {

	result := make(map[string]Record)

	data := z.serialize(true)

	var keys []string

	keys = []string{}

	keys = append(keys, data[feedNextFetch])

	keys = append(keys, data[feedID])

	result[feedIndexNextFetch] = Record{string(kvKeyEncode(keys...)): data[feedID]}

	keys = []string{}

	keys = append(keys, strings.ToLower(data[feedURL]))

	result[feedIndexURL] = Record{string(kvKeyEncode(keys...)): data[feedID]}

	return result
}

// GroupAllByURL groups collections of elements in Feeds by URL
func (z Feeds) GroupAllByURL() map[string]Feeds {
	result := make(map[string]Feeds)
	for _, feed := range z {
		a := result[feed.URL]
		a = append(a, feed)
		result[feed.URL] = a
	}
	return result
}
