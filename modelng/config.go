package modelng

const (
	entityConfig = "Config"
)

var (
	indexesConfig = []string{}
)

// Config defines an entry in the config bucket.
type Config struct {
	Name  string
	Value string
}
