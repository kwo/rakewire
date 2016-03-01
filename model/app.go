package model

import (
	"time"
)

const (
	logName  = "[model]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// application level variables
var (
	BuildHash string
	BuildTime string
	Version   string
	AppStart  time.Time
)