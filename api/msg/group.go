package msg

// Groups is a list of Group structs
type Groups []*Group

// Group defines a group of feeds
type Group struct {
	Name string `json:"name"`
}

// GroupListRequest defines a request to list groups for a specific user
type GroupListRequest struct{}

// GroupListResponse defines the response to a GroupListRequest
type GroupListResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Groups  Groups `json:"groups,omitempty"`
}
