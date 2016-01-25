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
	FeedLogIndexTime     = "Time"
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
	feedlogItemCount     = "ItemCount"
	feedlogNewItems      = "NewItems"
)

var (
	feedlogAllFields = []string{
		feedlogID, feedlogFeedID, feedlogDuration, feedlogResult, feedlogResultMessage, feedlogStartTime, feedlogURL, feedlogContentLength, feedlogContentType, feedlogETag, feedlogLastModified, feedlogStatusCode, feedlogUsesGzip, feedlogFlavor, feedlogGenerator, feedlogTitle, feedlogLastUpdated, feedlogItemCount, feedlogNewItems,
	}
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
	z.ItemCount = 0
	z.NewItems = 0

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *FeedLog) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[feedlogID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.FeedID != 0 {
		result[feedlogFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if flagNoZeroCheck || z.Duration != 0 {
		result[feedlogDuration] = z.Duration.String()
	}

	if flagNoZeroCheck || z.Result != "" {
		result[feedlogResult] = z.Result
	}

	if flagNoZeroCheck || z.ResultMessage != "" {
		result[feedlogResultMessage] = z.ResultMessage
	}

	if flagNoZeroCheck || !z.StartTime.IsZero() {
		result[feedlogStartTime] = z.StartTime.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.URL != "" {
		result[feedlogURL] = z.URL
	}

	if flagNoZeroCheck || z.ContentLength != 0 {
		result[feedlogContentLength] = fmt.Sprintf("%d", z.ContentLength)
	}

	if flagNoZeroCheck || z.ContentType != "" {
		result[feedlogContentType] = z.ContentType
	}

	if flagNoZeroCheck || z.ETag != "" {
		result[feedlogETag] = z.ETag
	}

	if flagNoZeroCheck || !z.LastModified.IsZero() {
		result[feedlogLastModified] = z.LastModified.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.StatusCode != 0 {
		result[feedlogStatusCode] = fmt.Sprintf("%d", z.StatusCode)
	}

	if flagNoZeroCheck || z.UsesGzip {
		result[feedlogUsesGzip] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.UsesGzip)
	}

	if flagNoZeroCheck || z.Flavor != "" {
		result[feedlogFlavor] = z.Flavor
	}

	if flagNoZeroCheck || z.Generator != "" {
		result[feedlogGenerator] = z.Generator
	}

	if flagNoZeroCheck || z.Title != "" {
		result[feedlogTitle] = z.Title
	}

	if flagNoZeroCheck || !z.LastUpdated.IsZero() {
		result[feedlogLastUpdated] = z.LastUpdated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.ItemCount != 0 {
		result[feedlogItemCount] = fmt.Sprintf("%d", z.ItemCount)
	}

	if flagNoZeroCheck || z.NewItems != 0 {
		result[feedlogNewItems] = fmt.Sprintf("%d", z.NewItems)
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *FeedLog) Deserialize(values map[string]string, flags ...bool) error {
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
	}(feedlogID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, feedlogID)
	}

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(feedlogFeedID, values, errors)

	if !(z.FeedID != 0) {
		missing = append(missing, feedlogFeedID)
	}

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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(feedlogUsesGzip, values, errors)

	z.Flavor = values[feedlogFlavor]

	z.Generator = values[feedlogGenerator]

	z.Title = values[feedlogTitle]

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
	}(feedlogLastUpdated, values, errors)

	z.ItemCount = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogItemCount, values, errors)

	z.NewItems = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(feedlogNewItems, values, errors)

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, feedlogAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *FeedLog) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[FeedLogIndexFeedTime] = []string{

		data[feedlogFeedID],

		data[feedlogStartTime],
	}

	result[FeedLogIndexTime] = []string{

		data[feedlogStartTime],

		data[feedlogID],
	}

	return result
}
