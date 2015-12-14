package model

import (
	"github.com/pborman/uuid"
	"sync"
)

const (
	fID   = "ID"
	empty = ""
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
