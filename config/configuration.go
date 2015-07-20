package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path"
	"rakewire.com/db"
	"rakewire.com/fetch"
	"rakewire.com/httpd"
	"rakewire.com/queryfeed"
)

const (
	configFileName = "config.yaml"
)

// Configuration object
type Configuration struct {
	Database  db.Configuration
	QueryFeed queryfeed.Configuration
	Fetcher   fetch.Configuration
	Httpd     httpd.Configuration
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

// GetConfig get the configurtion
func GetConfig() *Configuration {
	cfg := Configuration{}
	if err := cfg.LoadFromFile(GetConfigFileLocation()); err != nil {
		return nil
	}
	return &cfg
}

// GetConfigFileLocation get the location of the config file
func GetConfigFileLocation() string {
	if home := GetHomeDirectory(); home != "" {
		return path.Join(GetHomeDirectory(), ".rakewire", configFileName)
	}
	return configFileName
}

// GetHomeDirectory get the user home directory
func GetHomeDirectory() string {
	homeLocations := []string{"HOME", "HOMEPATH", "USERPROFILE"}
	for _, v := range homeLocations {
		x := os.Getenv(v)
		if x != "" {
			return x
		}
	}
	return ""
}
