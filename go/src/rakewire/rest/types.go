package rest

// User represents a user
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}
