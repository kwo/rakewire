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
	}
)
