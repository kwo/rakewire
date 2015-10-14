package logging

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

const (
	defaultPattern         = "%-8s %s %-8s %s\n"
	defaultTimestampFormat = "15:04:05"
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

// InternalFormatter for logging
type InternalFormatter struct {
	Pattern         string
	TimestampFormat string
	UseColor        bool
}

// NewInternalFormatter creates a new InternalFormatter
func NewInternalFormatter() *InternalFormatter {
	useColor := log.IsTerminal() && (runtime.GOOS != "windows")
	return &InternalFormatter{
		Pattern:         defaultPattern,
		TimestampFormat: defaultTimestampFormat,
		UseColor:        useColor,
	}
}

// Format log entry
func (z *InternalFormatter) Format(entry *log.Entry) ([]byte, error) {

	b := &bytes.Buffer{}

	level, levelColor := getLevelNameAndColor(entry.Level)
	timestamp := entry.Time.Format(z.TimestampFormat)
	category := entry.Data["category"]
	if category == nil || category == "" {
		category = "none"
	}
	message := strings.TrimSpace(strings.TrimSuffix(entry.Message, "\n"))

	if z.UseColor {
		level = fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, level)
	}

	fmt.Fprintf(b, z.Pattern, timestamp, level, category, message)

	return b.Bytes(), nil

}

func getLevelNameAndColor(level log.Level) (string, int) {

	var levelName string
	var levelColor int

	switch level {
	case log.DebugLevel:
		levelName = "DEBUG"
		levelColor = green
	case log.InfoLevel:
		levelName = "INFO "
		levelColor = cyan
	case log.WarnLevel:
		levelName = "WARN "
		levelColor = yellow
	case log.ErrorLevel:
		levelName = "ERROR"
		levelColor = red
	case log.FatalLevel:
		levelName = "FATAL"
		levelColor = red
	case log.PanicLevel:
		levelName = "PANIC"
		levelColor = red
	default:
		levelName = "NONE "
		levelColor = nocolor
	}

	return levelName, levelColor

}
