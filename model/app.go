package model

import (
	"time"
)

// app-level variables
var (
	Version   = "beta"
	BuildHash = ""
	BuildTime = time.Now().UTC().Format(time.RFC3339)
	AppStart  time.Time
)
