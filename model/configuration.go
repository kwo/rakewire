package model

import (
	"bytes"
	yaml "gopkg.in/yaml.v2"
	"io"
	"os"
)

// Configuration object
type Configuration struct {
	Httpd HttpdConfiguration
}

// HttpdConfiguration configuration
type HttpdConfiguration struct {
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
