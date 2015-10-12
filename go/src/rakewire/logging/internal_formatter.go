package logging

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
)

const (
	defaultPattern         = "%-8s %-8s - %s\n"
	defaultTimestampFormat = "15:04:05"
)

// InternalFormatter for logging
type InternalFormatter struct {
	Pattern         string
	TimestampFormat string
}

// NewInternalFormatter creates a new InternalFormatter
func NewInternalFormatter() *InternalFormatter {
	return &InternalFormatter{
		Pattern:         defaultPattern,
		TimestampFormat: defaultTimestampFormat,
	}
}

// Format log entry
func (z *InternalFormatter) Format(entry *log.Entry) ([]byte, error) {

	b := &bytes.Buffer{}

	category := entry.Data["category"]
	timestamp := entry.Time.Format(z.TimestampFormat)
	message := strings.TrimSpace(strings.TrimSuffix(entry.Message, "\n"))

	fmt.Fprintf(b, z.Pattern, category, timestamp, message)

	return b.Bytes(), nil

}
