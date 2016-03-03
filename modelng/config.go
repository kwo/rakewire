package modelng

const (
	entityConfig = "Config"
	idConfig     = "configuration"
)

var (
	indexesConfig = []string{}
)

// Config defines the application configurtion.
type Config struct {
	ID           string
	LoggingLevel string
	Sequences    sequences
}

type sequences struct {
	User  uint64
	Group uint64
}
