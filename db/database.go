package db

// Database interface
type Database interface {
	GetFeedByID(id string) (*Feed, error)
	GetFeeds() (*Feeds, error)
	SaveFeeds(*Feeds) (int, error)
}
