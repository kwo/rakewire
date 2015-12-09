package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// Configuration is a log writer
type Configuration struct {
	File    string
	Level   string
	NoColor bool
	level   *loglevel
	writer  io.Writer
}

type loglevel struct {
	name  string
	level int
	color int
}

// log level names
const (
	LogTrace = "TRACE"
	LogDebug = "DEBUG"
	LogInfo  = "INFO"
	LogWarn  = "WARN"
	LogError = "ERROR"
	logNone  = "NONE"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	cyan    = 36
	gray    = 37
)

var logLinePattern = regexp.MustCompile(`^(.*)\s+\[(\w+)\]\s+\[(\w+)\]\s+(.*)`)

var loglevels = map[string]*loglevel{
	LogTrace: &loglevel{
		name:  LogTrace,
		level: 1,
		color: gray,
	},
	LogDebug: &loglevel{
		name:  LogDebug,
		level: 2,
		color: green,
	},
	LogInfo: &loglevel{
		name:  LogInfo,
		level: 3,
		color: cyan,
	},
	LogWarn: &loglevel{
		name:  LogWarn,
		level: 4,
		color: yellow,
	},
	LogError: &loglevel{
		name:  LogError,
		level: 5,
		color: red,
	},
	logNone: &loglevel{
		name:  logNone,
		level: 6,
		color: nocolor,
	},
}

// Init the logging system
func (cfg *Configuration) Init() {

	if level := loglevels[cfg.Level]; level != nil {
		cfg.level = level
	} else {
		cfg.level = loglevels[LogInfo] // default
	}

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

	line := string(p)
	if matches := logLinePattern.FindStringSubmatch(line); len(matches) == 5 {

		linePrefix := matches[1]
		lineLevel := matches[2]
		lineCategory := matches[3]
		lineMessage := matches[4]

		loglevel := loglevels[lineLevel]
		if loglevel == nil {
			loglevel = loglevels[logNone] // default
		}

		if loglevel.level < cfg.level.level {
			return len(p), nil
		}

		if cfg.NoColor {
			line = fmt.Sprintf("%s %-5s %-5s  %s\n", linePrefix, loglevel.name, lineCategory, lineMessage)
		} else {
			line = fmt.Sprintf("%s \x1b[%dm%-5s\x1b[0m %-5s  %s\n", linePrefix, loglevel.color, loglevel.name, lineCategory, lineMessage)
		}

		return cfg.writer.Write([]byte(line))

	}

	return cfg.writer.Write(p)

}
