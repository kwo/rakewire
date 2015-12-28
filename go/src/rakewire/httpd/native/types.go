package native

// API top level struct
type API struct {
	prefix string
	db     Database
}

// Database defines the interface to the database
type Database interface {
}

// User represents a user
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}
