package main

import (
	"os"
	"path"
	m "rakewire.com/model"
)

func getConfig() *m.Configuration {
	cfg := m.Configuration{}
	if err := cfg.LoadFromFile(getConfigFileLocation()); err != nil {
		return nil
	}
	return &cfg
}

func getConfigFileLocation() string {
	if home := getHomeDirectory(); home != "" {
		return path.Join(getHomeDirectory(), ".rakewire", configFileName)
	}
	return configFileName
}

func getHomeDirectory() string {
	homeLocations := []string{"HOME", "HOMEPATH", "USERPROFILE"}
	for _, v := range homeLocations {
		x := os.Getenv(v)
		if x != "" {
			return x
		}
	}
	return ""
}
