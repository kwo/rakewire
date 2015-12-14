package model

import (
	"github.com/pborman/uuid"
	"sync"
	"time"
)

const (
	fID        = "ID"
	empty      = ""
	timeFormat = time.RFC3339Nano
)

// application level variables
var (
	BuildHash string
	BuildTime string
	Version   string
)

var uuidLock sync.Mutex

func getUUID() string {
	uuidLock.Lock()
	defer uuidLock.Unlock()
	return uuid.NewUUID().String()
}
