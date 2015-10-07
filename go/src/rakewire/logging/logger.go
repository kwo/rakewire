package logging

import (
	"fmt"
	"log"
	"os"
)

// Init the logging system
func Init(levelStr string) {
}

// New create a new internal logger
func New(category string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%-9.8s", category), log.Ltime)
}
