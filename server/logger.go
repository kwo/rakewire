package server

import (
	"github.com/codegangsta/negroni"
	"rakewire.com/logging"
)

// NewInternalLogger returns a new Logger instance
func NewInternalLogger() *negroni.Logger {
	return &negroni.Logger{logging.New("server")}
}
