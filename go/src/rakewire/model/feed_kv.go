package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"sort"
	"strconv"
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
func (z *Feed) getID() uint64 {
	return z.ID
}

// setID sets the primary key of the object.
func (z *Feed) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Feed) clear() {
	z.ID = 0
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

	if flagNoZeroCheck || z.ID != 0 {
		result[feedID] = fmt.Sprintf("%05d", z.ID)
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
		result[feedLastModified] = z.LastModified.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || !z.LastUpdated.IsZero() {
		result[feedLastUpdated] = z.LastUpdated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || !z.NextFetch.IsZero() {
		result[feedNextFetch] = z.NextFetch.UTC().Format(time.RFC3339)
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
		result[feedStatusSince] = z.StatusSince.UTC().Format(time.RFC3339)
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

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(feedID, values, errors)

	if !(z.ID != 0) {
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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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

func newFeedID(tx Transaction) (string, error) {
	return kvNextID(feedEntity, tx)
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
