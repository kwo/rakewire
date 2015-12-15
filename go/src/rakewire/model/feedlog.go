package model

import (
	"time"
)

// FeedLog represents an attempted HTTP request to a feed
type FeedLog struct {
	ID            uint64        `json:"id"`
	FeedID        string        `json:"feedId"`
	Duration      time.Duration `json:"duration"`
	Result        string        `json:"result"`
	ResultMessage string        `json:"resultMessage"`
	StartTime     time.Time     `json:"startTime"`
	URL           string        `json:"url"`
	ContentLength int           `json:"contentLength"`
	ContentType   string        `json:"contentType"`
	ETag          string        `json:"etag"`
	LastModified  time.Time     `json:"lastModified"`
	StatusCode    int           `json:"statusCode"`
	UsesGzip      bool          `json:"gzip"`
	Flavor        string        `json:"flavor"`
	Generator     string        `json:"generator"`
	Title         string        `json:"title"`
	LastUpdated   time.Time     `json:"lastUpdated"`
	EntryCount    int           `json:"entryCount"`
	NewEntries    int           `json:"newEntries"`
}

// index constants
const (
	FeedLogEntity        = "FeedLog"
	FeedLogIndexFeedTime = "FeedTime"
)

// FetchResults
const (
	FetchResultOK          = "OK"
	FetchResultRedirect    = "MV" // message contains old URL -> new URL
	FetchResultClientError = "EC" // message contains error text
	FetchResultServerError = "ES" // check http status code
	FetchResultFeedError   = "FP" // cannot parse feed
)

const (
	flID            = "ID"
	flFeedID        = "FeedID"
	flDuration      = "Duration"
	flResult        = "Result"
	flResultMessage = "ResultMessage"
	flStartTime     = "StartTime"
	flURL           = "URL"
	flContentLength = "ContentLength"
	flContentType   = "ContentType"
	flETag          = "ETag"
	flLastModified  = "LastModified"
	flStatusCode    = "StatusCode"
	flUsesGzip      = "UsesGzip"
	flFlavor        = "Flavor"
	flGenerator     = "Generator"
	flTitle         = "Title"
	flLastUpdated   = "LastUpdated"
	flEntryCount    = "EntryCount"
	flNewEntries    = "NewEntries"
)

// NewFeedLog instantiates a new FeedLog with the required fields set.
func NewFeedLog(feedID string) *FeedLog {
	return &FeedLog{
		FeedID: feedID,
	}
}

// GetName return the name of the entity.
func (z *FeedLog) GetName() string {
	return FeedLogEntity
}

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
	z.EntryCount = 0
	z.NewEntries = 0
}

// Serialize serializes an object to a list of key-values.
func (z *FeedLog) Serialize() map[string]string {
	result := make(map[string]string)
	setUint(z.ID, flID, result)
	setString(z.FeedID, flFeedID, result)
	setDuration(z.Duration, flDuration, result)
	setString(z.Result, flResult, result)
	setString(z.ResultMessage, flResultMessage, result)
	setTime(z.StartTime, flStartTime, result)
	setString(z.URL, flURL, result)
	setInt(z.ContentLength, flContentLength, result)
	setString(z.ContentType, flContentType, result)
	setString(z.ETag, flETag, result)
	setTime(z.LastModified, flLastModified, result)
	setInt(z.StatusCode, flStatusCode, result)
	setBool(z.UsesGzip, flUsesGzip, result)
	setString(z.Flavor, flFlavor, result)
	setString(z.Generator, flGenerator, result)
	setString(z.Title, flTitle, result)
	setTime(z.LastUpdated, flLastUpdated, result)
	setInt(z.EntryCount, flEntryCount, result)
	setInt(z.NewEntries, flNewEntries, result)
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *FeedLog) Deserialize(values map[string]string) error {
	var errors []error
	z.ID = getUint(flID, values, errors)
	z.FeedID = getString(flFeedID, values, errors)
	z.Duration = getDuration(flDuration, values, errors)
	z.Result = getString(flResult, values, errors)
	z.ResultMessage = getString(flResultMessage, values, errors)
	z.StartTime = getTime(flStartTime, values, errors)
	z.URL = getString(flURL, values, errors)
	z.ContentLength = getInt(flContentLength, values, errors)
	z.ContentType = getString(flContentType, values, errors)
	z.ETag = getString(flETag, values, errors)
	z.LastModified = getTime(flLastModified, values, errors)
	z.StatusCode = getInt(flStatusCode, values, errors)
	z.UsesGzip = getBool(flUsesGzip, values, errors)
	z.Flavor = getString(flFlavor, values, errors)
	z.Generator = getString(flGenerator, values, errors)
	z.Title = getString(flTitle, values, errors)
	z.LastUpdated = getTime(flLastUpdated, values, errors)
	z.EntryCount = getInt(flEntryCount, values, errors)
	z.NewEntries = getInt(flNewEntries, values, errors)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil

}

// IndexKeys returns the keys of all indexes for this object.
func (z *FeedLog) IndexKeys() map[string][]string {
	data := z.Serialize()
	result := make(map[string][]string)
	result[FeedLogIndexFeedTime] = []string{data[flFeedID], data[flStartTime]}
	return result
}
