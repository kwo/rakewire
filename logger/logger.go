package logger

import (
	"fmt"
	"time"
)

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
)

// Package variables
var (
	DateFormat = "2006-02-01 15:04:05"
	DebugMode  = false
)

// New creates a new Logger
func New(name string) *Logger {
	return &Logger{Name: name}
}

// Logger defines the logger
type Logger struct {
	Silent bool
	Name   string
}

// Debugf writes a debug message to the console.
func (z *Logger) Debugf(format string, a ...interface{}) {
	if !z.Silent && DebugMode {
		z.logf(levelDebug, format, a)
	}
}

// Infof writes an info message to the console.
func (z *Logger) Infof(format string, a ...interface{}) {
	if !z.Silent {
		z.logf(levelInfo, format, a)
	}
}

func (z *Logger) logf(level, format0 string, a []interface{}) {

	var format = "%s %s: " + format0 + "\n"
	var args []interface{}
	args = append(args, time.Now().Format(DateFormat))
	args = append(args, z.Name)
	for _, arg := range a {
		args = append(args, arg)
	}
	fmt.Printf(format, args...)

}
