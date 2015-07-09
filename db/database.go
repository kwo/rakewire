package db

// Database interface
type Database interface {
	// return feeds keyed by ID
	GetFeeds() (map[string]*FeedInfo, error)
	SaveFeeds([]*FeedInfo) (int, error)
}
