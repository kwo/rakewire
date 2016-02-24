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
	transmissionEntity        = "Transmission"
	transmissionIndexFeedTime = "FeedTime"
	transmissionIndexTime     = "Time"
)

const (
	transmissionID            = "ID"
	transmissionFeedID        = "FeedID"
	transmissionDuration      = "Duration"
	transmissionResult        = "Result"
	transmissionResultMessage = "ResultMessage"
	transmissionStartTime     = "StartTime"
	transmissionURL           = "URL"
	transmissionContentLength = "ContentLength"
	transmissionContentType   = "ContentType"
	transmissionETag          = "ETag"
	transmissionLastModified  = "LastModified"
	transmissionStatusCode    = "StatusCode"
	transmissionUsesGzip      = "UsesGzip"
	transmissionFlavor        = "Flavor"
	transmissionGenerator     = "Generator"
	transmissionTitle         = "Title"
	transmissionLastUpdated   = "LastUpdated"
	transmissionItemCount     = "ItemCount"
	transmissionNewItems      = "NewItems"
)

var (
	transmissionAllFields = []string{
		transmissionID, transmissionFeedID, transmissionDuration, transmissionResult, transmissionResultMessage, transmissionStartTime, transmissionURL, transmissionContentLength, transmissionContentType, transmissionETag, transmissionLastModified, transmissionStatusCode, transmissionUsesGzip, transmissionFlavor, transmissionGenerator, transmissionTitle, transmissionLastUpdated, transmissionItemCount, transmissionNewItems,
	}
	transmissionAllIndexes = []string{
		transmissionIndexFeedTime, transmissionIndexTime,
	}
)

// Transmissions is a collection of Transmission elements
type Transmissions []*Transmission

func (z Transmissions) Len() int      { return len(z) }
func (z Transmissions) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Transmissions) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Transmissions) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Transmissions) First() *Transmission { return z[0] }

// Reverse reverses the order of the collection
func (z Transmissions) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Transmission) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Transmission) clear() {
	z.ID = empty
	z.FeedID = empty
	z.Duration = 0
	z.Result = empty
	z.ResultMessage = empty
	z.StartTime = time.Time{}
	z.URL = empty
	z.ContentLength = 0
	z.ContentType = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.StatusCode = 0
	z.UsesGzip = false
	z.Flavor = empty
	z.Generator = empty
	z.Title = empty
	z.LastUpdated = time.Time{}
	z.ItemCount = 0
	z.NewItems = 0

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Transmission) serialize(flags ...bool) Record {

	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(Record)

	if flagNoZeroCheck || z.ID != empty {
		result[transmissionID] = z.ID
	}

	if flagNoZeroCheck || z.FeedID != empty {
		result[transmissionFeedID] = z.FeedID
	}

	if flagNoZeroCheck || z.Duration != 0 {
		result[transmissionDuration] = z.Duration.String()
	}

	if flagNoZeroCheck || z.Result != empty {
		result[transmissionResult] = z.Result
	}

	if flagNoZeroCheck || z.ResultMessage != empty {
		result[transmissionResultMessage] = z.ResultMessage
	}

	if flagNoZeroCheck || !z.StartTime.IsZero() {
		result[transmissionStartTime] = z.StartTime.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || z.URL != empty {
		result[transmissionURL] = z.URL
	}

	if flagNoZeroCheck || z.ContentLength != 0 {
		result[transmissionContentLength] = fmt.Sprintf("%d", z.ContentLength)
	}

	if flagNoZeroCheck || z.ContentType != empty {
		result[transmissionContentType] = z.ContentType
	}

	if flagNoZeroCheck || z.ETag != empty {
		result[transmissionETag] = z.ETag
	}

	if flagNoZeroCheck || !z.LastModified.IsZero() {
		result[transmissionLastModified] = z.LastModified.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || z.StatusCode != 0 {
		result[transmissionStatusCode] = fmt.Sprintf("%d", z.StatusCode)
	}

	if flagNoZeroCheck || z.UsesGzip {
		result[transmissionUsesGzip] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.UsesGzip)
	}

	if flagNoZeroCheck || z.Flavor != empty {
		result[transmissionFlavor] = z.Flavor
	}

	if flagNoZeroCheck || z.Generator != empty {
		result[transmissionGenerator] = z.Generator
	}

	if flagNoZeroCheck || z.Title != empty {
		result[transmissionTitle] = z.Title
	}

	if flagNoZeroCheck || !z.LastUpdated.IsZero() {
		result[transmissionLastUpdated] = z.LastUpdated.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || z.ItemCount != 0 {
		result[transmissionItemCount] = fmt.Sprintf("%d", z.ItemCount)
	}

	if flagNoZeroCheck || z.NewItems != 0 {
		result[transmissionNewItems] = fmt.Sprintf("%d", z.NewItems)
	}

	return result

}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Transmission) deserialize(values Record, flags ...bool) error {

	flagUnknownCheck := len(flags) > 0 && flags[0]
	z.clear()

	var errors []error
	var missing []string
	var unknown []string

	z.ID = values[transmissionID]
	if !(z.ID != empty) {
		missing = append(missing, transmissionID)
	}

	z.FeedID = values[transmissionFeedID]
	if !(z.FeedID != empty) {
		missing = append(missing, transmissionFeedID)
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
	}(transmissionDuration, values, errors)

	z.Result = values[transmissionResult]
	z.ResultMessage = values[transmissionResultMessage]
	z.StartTime = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(transmissionStartTime, values, errors)

	z.URL = values[transmissionURL]
	z.ContentLength = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(transmissionContentLength, values, errors)

	z.ContentType = values[transmissionContentType]
	z.ETag = values[transmissionETag]
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
	}(transmissionLastModified, values, errors)

	z.StatusCode = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(transmissionStatusCode, values, errors)

	z.UsesGzip = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(transmissionUsesGzip, values, errors)

	z.Flavor = values[transmissionFlavor]
	z.Generator = values[transmissionGenerator]
	z.Title = values[transmissionTitle]
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
	}(transmissionLastUpdated, values, errors)

	z.ItemCount = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(transmissionItemCount, values, errors)

	z.NewItems = func(fieldName string, values map[string]string, errors []error) int {
		result, err := strconv.ParseInt(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return int(result)
	}(transmissionNewItems, values, errors)

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, transmissionAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}

	return newDeserializationError(transmissionEntity, errors, missing, unknown)

}

// serializeIndexes returns all index records
func (z *Transmission) serializeIndexes() map[string]Record {

	result := make(map[string]Record)
	data := z.serialize(true)
	var keys []string

	keys = []string{}
	keys = append(keys, data[transmissionFeedID])
	keys = append(keys, data[transmissionStartTime])
	result[transmissionIndexFeedTime] = Record{string(kvKeyEncode(keys...)): data[transmissionID]}

	keys = []string{}
	keys = append(keys, data[transmissionStartTime])
	keys = append(keys, data[transmissionID])
	result[transmissionIndexTime] = Record{string(kvKeyEncode(keys...)): data[transmissionID]}

	return result

}
