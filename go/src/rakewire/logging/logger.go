package logging

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

// Configuration for logging
type Configuration struct {
	Level           string
	Pattern         string
	TimestampFormat string
	NoColor         bool
}

// Init the logging system
func Init(cfg *Configuration) {

	log.SetOutput(os.Stdout)

	formatter := NewInternalFormatter()

	if cfg.Pattern != "" {
		formatter.Pattern = cfg.Pattern
	}

	if cfg.TimestampFormat != "" {
		formatter.TimestampFormat = cfg.TimestampFormat
	}

	if cfg.NoColor {
		formatter.UseColor = false
	}

	log.SetFormatter(formatter)

	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

}

// New create a new internal logger
func New(category string) *log.Entry {
	return log.WithField("category", category)
}
