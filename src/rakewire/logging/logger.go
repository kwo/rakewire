package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
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

// Null create a new internal logger
func Null(category string) *log.Logger {
	tmp, _ := os.Open(os.DevNull)
	return log.New(tmp, fmt.Sprintf("%-9.8s", category), log.Ltime)
}

// NewFromType create a new internal logger
func NewFromType(a interface{}) *log.Logger {
	return New(fmt.Sprintf("%-9.8s", path.Base(reflect.TypeOf(a).PkgPath())))
}
