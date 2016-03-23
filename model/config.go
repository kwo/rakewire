package model

const (
	entityConfig = "Config"
	idConfig     = "configuration"
)

var (
	indexesConfig = []string{}
)

// Config defines the application configurtion.
type Config struct {
	ID        string
	Sequences sequences
	Log       logConfig
}

type logConfig struct {
	Level string
}

type sequences struct {
	User         uint64
	Feed         uint64
	Item         uint64
	Group        uint64
	Transmission uint64
}
