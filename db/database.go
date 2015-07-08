package db

// Database interface
type Database interface {
	init(cfg string) error
	destroy() error
	// return feeds keyed by ID
	getFeeds() (map[string]*FeedInfo, error)
}
