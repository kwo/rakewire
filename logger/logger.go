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
		z.logf(levelDebug, format, a...)
	}
}

// Infof writes an info message to the console.
func (z *Logger) Infof(format string, a ...interface{}) {
	if !z.Silent {
		z.logf(levelInfo, format, a...)
	}
}

func (z *Logger) logf(level, format string, a ...interface{}) {
	dateStr := time.Now().Format(DateFormat)
	if DebugMode {
		fmt.Printf("%s %-5s %-10s - ", dateStr, level, z.Name)
	} else {
		fmt.Printf("%s %-10s - ", dateStr, z.Name)
	}
	fmt.Printf(format, a...)
	fmt.Println()
}
