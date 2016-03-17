package model

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

// ByID groups elements in the Groups collection by ID
func (z Groups) ByID() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.ID] = group
	}
	return result
}

// ByName groups elements in the Groups collection by Name
func (z Groups) ByName() map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range z {
		result[group.Name] = group
	}
	return result
}

// Group defines an item status for a user
type Group struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
}
