package server

import (
	"github.com/codegangsta/negroni"
	"log"
	"os"
)

// NewInternalLogger returns a new Logger instance
func NewInternalLogger() *negroni.Logger {
	return &negroni.Logger{log.New(os.Stdout, "[rakewire] ", 0)}
}
