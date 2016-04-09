package model

//go:generate go run ${GOPATH}/src/rakewire/tools/buildinfo/buildinfo.go

import (
	"time"
)

// Version number of the app
const Version = "1.10.0"

// AppStart marks the time the application was started.
var AppStart time.Time
