package logging

import (
	"fmt"
	"log"
	"os"
)

// TODO: replace logrus with logutils

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

	switch cfg.File {
	case "stderr":
		log.SetOutput(os.Stderr)
	case "", "stdout":
		log.SetOutput(os.Stdout)
	default:
		f, err := os.OpenFile(cfg.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err == nil {
			log.SetOutput(f)
		} else {
			fmt.Printf("Reverting to stdout, cannot open log file: %s", err.Error())
			log.SetOutput(os.Stdout)
		}
	}

}
