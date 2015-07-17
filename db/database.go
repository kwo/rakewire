package db

// Database interface
type Database interface {
	GetFeedByID(id string) (*Feed, error)
	GetFeedByURL(url string) (*Feed, error)
	GetFeeds() (*Feeds, error)
	SaveFeeds(*Feeds) (int, error)
	Repair() error
}

// Configuration configuration
type Configuration struct {
	Location string
}
