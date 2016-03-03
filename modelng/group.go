package modelng

const (
	entityGroup        = "Group"
	indexGroupUserName = "UserName"
)

var (
	indexesGroup = []string{
		indexGroupUserName,
	}
)

// Groups is a collection of Group elements
type Groups []*Group

// Group defines an item status for a user
type Group struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
	Name   string `json:"name"`
}
