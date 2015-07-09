package logging

import (
	"fmt"
	"log"
	"os"
)

var (
	// Output file
	output = os.Stdout
)

// Linefeed skip a line
func Linefeed() {
	output.WriteString("\n")
}

// New create a new internal logger
func New(category string) *log.Logger {
	return log.New(output, fmt.Sprintf("%-9.8s", category), log.Ltime)
}
