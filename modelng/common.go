package modelng

const (
	bucketData  = "Data"
	bucketIndex = "Index"
	chSep       = "|"
	empty       = ""
	fmtTime     = "20060102150405Z0700"
	fmtUint     = "%010d"
)

var (
	allEntities = map[string][]string{
		entityConfig: indexesConfig,
		entityUser:   indexesUser,
		entityEntry:  indexesEntry,
	}
)

func getObject(entityName string) Object {
	switch entityName {
	case entityConfig:
		return &Config{}
	case entityEntry:
		return &Entry{}
	case entityUser:
		return &User{}
	}
	return nil
}
