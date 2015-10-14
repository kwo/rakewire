package logging

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

// Configuration for logging
type Configuration struct {
	File            string
	Level           string
	Pattern         string
	TimestampFormat string
	NoColor         bool
}

// Init the logging system
func Init(cfg *Configuration) {

	var noColor bool
	var logFile *os.File
	switch cfg.File {
	case "stderr":
		log.SetOutput(os.Stderr)
		logFile = os.Stderr
	case "", "stdout":
		log.SetOutput(os.Stdout)
		logFile = os.Stdout
	default:
		f, err := os.OpenFile(cfg.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err == nil {
			log.SetOutput(f)
			logFile = f
			noColor = true
		} else {
			fmt.Printf("Reverting to stdout, cannot open log file: %s", err.Error())
			log.SetOutput(os.Stdout)
			logFile = os.Stdout
		}
	}

	formatter := NewInternalFormatter()

	if cfg.Pattern != "" {
		formatter.Pattern = cfg.Pattern
	}

	if cfg.TimestampFormat != "" {
		formatter.TimestampFormat = cfg.TimestampFormat
	}

	if cfg.NoColor || noColor {
		formatter.UseColor = false
	}

	log.SetFormatter(formatter)

	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	log.Infof("Logging at %s level to %s", level.String(), logFile.Name())

}

// New create a new internal logger
func New(category string) *log.Entry {
	return log.WithField("category", category)
}
