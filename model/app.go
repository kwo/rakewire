package model

import (
	"time"
)

// Version number of the app
const Version = "1.10.0"

// AppStart marks the time the application was started.
var AppStart time.Time

// Last Commit variables
var (
	CommitHash = "<COMMITHASH>"
	CommitTime = "<COMMITTERDATEISO8601>"
)
