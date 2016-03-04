package modelng

const (
	bucketData  = "Data"
	bucketIndex = "Index"
	chMax       = "~"
	chSep       = "|"
	empty       = ""
	fmtTime     = "20060102150405Z0700"
	fmtUint     = "%010d"
)

var (
	allEntities = map[string][]string{
		entityConfig:       indexesConfig,
		entityEntry:        indexesEntry,
		entityFeed:         indexesFeed,
		entityGroup:        indexesGroup,
		entityItem:         indexesItem,
		entityTransmission: indexesTransmission,
		entityUser:         indexesUser,
	}
)

func getObject(entityName string) Object {
	switch entityName {
	case entityConfig:
		return &Config{}
	case entityEntry:
		return &Entry{}
	case entityFeed:
		return &Feed{}
	case entityGroup:
		return &Group{}
	case entityItem:
		return &Item{}
	case entityTransmission:
		return &Transmission{}
	case entityUser:
		return &User{}
	}
	return nil
}
