package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// index names
const (
	FeedEntity         = "Feed"
	FeedIndexNextFetch = "NextFetch"
	FeedIndexURL       = "URL"
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

// GetID return the primary key of the object.
func (z *Feed) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Feed) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Feed) Clear() {
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
func (z *Feed) Serialize(flags ...bool) map[string]string {
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
func (z *Feed) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(feedID, values, errors)

	z.URL = values[feedURL]

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

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Feed) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[FeedIndexNextFetch] = []string{

		data[feedNextFetch],

		data[feedID],
	}

	result[FeedIndexURL] = []string{

		strings.ToLower(data[feedURL]),
	}

	return result
}
