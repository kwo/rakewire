package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"strconv"
	"time"
)

// index names
const (
	FeedLogEntity        = "FeedLog"
	FeedLogIndexFeedTime = "FeedTime"
)

const (
	feedlogID            = "ID"
	feedlogFeedID        = "FeedID"
	feedlogDuration      = "Duration"
	feedlogResult        = "Result"
	feedlogResultMessage = "ResultMessage"
	feedlogStartTime     = "StartTime"
	feedlogURL           = "URL"
	feedlogContentLength = "ContentLength"
	feedlogContentType   = "ContentType"
	feedlogETag          = "ETag"
	feedlogLastModified  = "LastModified"
	feedlogStatusCode    = "StatusCode"
	feedlogUsesGzip      = "UsesGzip"
	feedlogFlavor        = "Flavor"
	feedlogGenerator     = "Generator"
	feedlogTitle         = "Title"
	feedlogLastUpdated   = "LastUpdated"
	feedlogEntryCount    = "EntryCount"
	feedlogNewEntries    = "NewEntries"
)

// GetID return the primary key of the object.
func (z *FeedLog) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *FeedLog) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *FeedLog) Clear() {
	z.ID = 0
	z.FeedID = 0
	z.Duration = 0
	z.Result = ""
	z.ResultMessage = ""
	z.StartTime = time.Time{}
	z.URL = ""
	z.ContentLength = 0
	z.ContentType = ""
	z.ETag = ""
	z.LastModified = time.Time{}
	z.StatusCode = 0
	z.UsesGzip = false
	z.Flavor = ""
	z.Generator = ""
	z.Title = ""
	z.LastUpdated = time.Time{}
	z.EntryCount = 0
	z.NewEntries = 0

}

// Serialize serializes an object to a list of key-values.
func (z *FeedLog) Serialize() map[string]string {
	result := make(map[string]string)

	if z.ID != 0 {
		result[feedlogID] = fmt.Sprintf("%05d", z.ID)
	}

	if z.FeedID != 0 {
		result[feedlogFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if z.Duration != 0 {
		result[feedlogDuration] = z.Duration.String()
	}

	if z.Result != "" {
		result[feedlogResult] = z.Result
	}

	if z.ResultMessage != "" {
		result[feedlogResultMessage] = z.ResultMessage
	}

	if !z.StartTime.IsZero() {
		result[feedlogStartTime] = z.StartTime.Format(time.RFC3339Nano)
	}

	if z.URL != "" {
		result[feedlogURL] = z.URL
	}

	if z.ContentLength != 0 {
		result[feedlogContentLength] = fmt.Sprintf("%d", z.ContentLength)
	}

	if z.ContentType != "" {
		result[feedlogContentType] = z.ContentType
	}

	if z.ETag != "" {
		result[feedlogETag] = z.ETag
	}

	if !z.LastModified.IsZero() {
		result[feedlogLastModified] = z.LastModified.Format(time.RFC3339Nano)
	}

	if z.StatusCode != 0 {
		result[feedlogStatusCode] = fmt.Sprintf("%d", z.StatusCode)
	}

	if z.UsesGzip {
		result[feedlogUsesGzip] = fmt.Sprintf("%t", z.UsesGzip)
	}

	if z.Flavor != "" {
		result[feedlogFlavor] = z.Flavor
	}

	if z.Generator != "" {
		result[feedlogGenerator] = z.Generator
	}

	if z.Title != "" {
		result[feedlogTitle] = z.Title
	}

	if !z.LastUpdated.IsZero() {
		result[feedlogLastUpdated] = z.LastUpdated.Format(time.RFC3339Nano)
	}

	if z.EntryCount != 0 {
		result[feedlogEntryCount] = fmt.Sprintf("%d", z.EntryCount)
	}

	if z.NewEntries != 0 {
		result[feedlogNewEntries] = fmt.Sprintf("%d", z.NewEntries)
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *FeedLog) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(feedlogID, values, errors)

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(feedlogFeedID, values, errors)

	z.Duration = func(fieldName string, values map[string]string, errors []error) time.Duration {
		var result time.Duration
		if value, ok := values[fieldName]; ok {
			t, err := time.ParseDuration(value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(feedlogDuration, values, errors)

	z.Result = values[feedlogResult]

	z.ResultMessage = values[feedlogResultMessage]

	z.StartTime = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(feedlogStartTime, values, errors)

	z.URL = values[feedlogURL]

	z.ContentLength = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogContentLength, values, errors)

	z.ContentType = values[feedlogContentType]

	z.ETag = values[feedlogETag]

	z.LastModified = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(feedlogLastModified, values, errors)

	z.StatusCode = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogStatusCode, values, errors)

	z.UsesGzip = func(fieldName string, values map[string]string, errors []error) bool {
		result, err := strconv.ParseBool(values[fieldName])
		if err != nil {
			errors = append(errors, err)
			return false
		}
		return result
	}(feedlogUsesGzip, values, errors)

	z.Flavor = values[feedlogFlavor]

	z.Generator = values[feedlogGenerator]

	z.Title = values[feedlogTitle]

	z.LastUpdated = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(feedlogLastUpdated, values, errors)

	z.EntryCount = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogEntryCount, values, errors)

	z.NewEntries = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogNewEntries, values, errors)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *FeedLog) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize()

	result[FeedLogIndexFeedTime] = []string{

		data[feedlogFeedID],

		data[feedlogStartTime],
	}

	return result
}
