package model

const (
	// VERSION application version
	VERSION = "0.1.0"
)

// Identifiable marks structs with an ID
type Identifiable interface {
	GetID() string
}
