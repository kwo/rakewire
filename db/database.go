package db

// Database interface
type Database interface {
	GetFeeds() (*Feeds, error)
	SaveFeeds(*Feeds) (int, error)
}
