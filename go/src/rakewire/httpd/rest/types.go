package rest

import (
	"rakewire/db"
)

// API top level struct
type API struct {
	prefix string
	db     db.Database
}

// User represents a user
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}
