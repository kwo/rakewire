package model

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Configuration object
type Configuration struct {
	Database DatabaseConfiguration
	Fetcher  FetcherConfiguration
	Httpd    HttpdConfiguration
}

// DatabaseConfiguration configuration
type DatabaseConfiguration struct {
	Location string
}

// FetcherConfiguration configuration
type FetcherConfiguration struct {
	Fetchers           int
	RequestBuffer      int
	HTTPTimeoutSeconds int
	IdleTimeoutSeconds int
}

// HttpdConfiguration configuration
type HttpdConfiguration struct {
	Address   string
	Port      int
	WebAppDir string
}

// Load configuration from JSON
func (z *Configuration) Load(data []byte) error {
	return yaml.Unmarshal(data, z)
}

// LoadFromReader configuration from JSON
func (z *Configuration) LoadFromReader(r io.ReadCloser) error {
	var b bytes.Buffer
	b.ReadFrom(r)
	r.Close()
	return z.Load(b.Bytes())
}

// LoadFromFile configuration from JSON
func (z *Configuration) LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return z.LoadFromReader(f)
}
