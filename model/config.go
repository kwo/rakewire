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
	Fetch     fetchConfig
	Log       logConfig
}

type fetchConfig struct {
	Timeout string
	Workers int
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

// GetInt returns the given value if nonzero otherwise the default value
func (cfg *Config) GetInt(value int, defaultValue int) int {
	if value != 0 {
		return value
	}
	return defaultValue
}

// GetStr returns the given value if not empty otherwise the default value
func (cfg *Config) GetStr(value string, defaultValue string) string {
	if value != empty {
		return value
	}
	return defaultValue
}
