package logging

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// Configuration is a log writer
type Configuration struct {
	File      string
	Level     string
	Levels    []string
	writer    io.Writer
	badLevels map[string]struct{}
}

// Init the logging system
func (cfg *Configuration) Init() {

	cfg.Levels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

	badLevels := make(map[string]struct{})
	for _, level := range cfg.Levels {
		if level == cfg.Level {
			break
		}
		badLevels[level] = struct{}{}
	}
	cfg.badLevels = badLevels

	switch cfg.File {
	case "stderr":
		cfg.writer = os.Stderr
	case "", "stdout":
		cfg.writer = os.Stdout
	default:
		f, err := os.OpenFile(cfg.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err == nil {
			cfg.writer = f
		} else {
			fmt.Printf("Reverting to stdout, cannot open log file: %s", err.Error())
			cfg.writer = os.Stdout
		}
	}

	log.SetOutput(cfg)

}

func (cfg *Configuration) Write(p []byte) (n int, err error) {

	if !cfg.check(p) {
		return len(p), nil
	}

	return cfg.writer.Write(p)
}

func (cfg *Configuration) check(line []byte) bool {

	// Check for a log level
	var level string
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = string(line[x+1 : x+y])
		}
	}

	_, ok := cfg.badLevels[level]
	return !ok

}
