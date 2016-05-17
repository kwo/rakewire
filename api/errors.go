package api

import (
	"errors"
)

// package level senteniel errors
var (
	errEscape       = errors.New("Escape Error")
	ErrEmptyRequest = errors.New("Empty Request")
)
