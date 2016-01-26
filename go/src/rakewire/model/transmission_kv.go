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
)

// Transmissions is a collection of Transmission elements
type Transmissions []*Transmission

func (z Transmissions) Len() int      { return len(z) }
func (z Transmissions) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Transmissions) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// First returns the first element in the collection
func (z Transmissions) First() *Transmission { return z[0] }

// Reverse reverses the order of the collection
func (z Transmissions) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// GetID return the primary key of the object.
func (z *Transmission) getID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Transmission) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Transmission) clear() {
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
func (z *Transmission) serialize(flags ...bool) Record {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[transmissionID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.FeedID != 0 {
		result[transmissionFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if flagNoZeroCheck || z.Duration != 0 {
		result[transmissionDuration] = z.Duration.String()
	}

	if flagNoZeroCheck || z.Result != "" {
		result[transmissionResult] = z.Result
	}

	if flagNoZeroCheck || z.ResultMessage != "" {
		result[transmissionResultMessage] = z.ResultMessage
	}

	if flagNoZeroCheck || !z.StartTime.IsZero() {
		result[transmissionStartTime] = z.StartTime.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.URL != "" {
		result[transmissionURL] = z.URL
	}

	if flagNoZeroCheck || z.ContentLength != 0 {
		result[transmissionContentLength] = fmt.Sprintf("%d", z.ContentLength)
	}

	if flagNoZeroCheck || z.ContentType != "" {
		result[transmissionContentType] = z.ContentType
	}

	if flagNoZeroCheck || z.ETag != "" {
		result[transmissionETag] = z.ETag
	}

	if flagNoZeroCheck || !z.LastModified.IsZero() {
		result[transmissionLastModified] = z.LastModified.UTC().Format(time.RFC3339)
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

	if flagNoZeroCheck || z.Flavor != "" {
		result[transmissionFlavor] = z.Flavor
	}

	if flagNoZeroCheck || z.Generator != "" {
		result[transmissionGenerator] = z.Generator
	}

	if flagNoZeroCheck || z.Title != "" {
		result[transmissionTitle] = z.Title
	}

	if flagNoZeroCheck || !z.LastUpdated.IsZero() {
		result[transmissionLastUpdated] = z.LastUpdated.UTC().Format(time.RFC3339)
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
	}(transmissionID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, transmissionID)
	}

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(transmissionFeedID, values, errors)

	if !(z.FeedID != 0) {
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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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
			t, err := time.Parse(time.RFC3339, value)
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

// IndexKeys returns the keys of all indexes for this object.
func (z *Transmission) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[transmissionIndexFeedTime] = []string{

		data[transmissionFeedID],

		data[transmissionStartTime],
	}

	result[transmissionIndexTime] = []string{

		data[transmissionStartTime],

		data[transmissionID],
	}

	return result
}
