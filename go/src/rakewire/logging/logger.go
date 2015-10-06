package logging

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

// Init the logging system
func Init(levelStr string) {
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

// New create a new internal logger
func New(category string) *log.Entry {
	return log.WithFields(log.Fields{
		"category": category,
	})
}
