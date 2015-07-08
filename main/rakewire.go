package main

import (
	"fmt"
	"os"
	"path"
	m "rakewire.com/model"
	"rakewire.com/server"
)

const (
	configFileName = "config.yaml"
)

func main() {

	cfg := getConfig()
	if cfg == nil {
		fmt.Printf("Abort! No config file found at %s\n", getConfigFileLocation())
		os.Exit(1)
		return
	}

	server.Serve(cfg.Httpd)

}

func getConfig() *m.Configuration {
	cfg := m.Configuration{}
	if err := cfg.LoadFromFile(getConfigFileLocation()); err != nil {
		return nil
	}
	return &cfg
}

func getConfigFileLocation() string {
	if home := getHomeDirectory(); home != "" {
		return path.Join(getHomeDirectory(), ".config", "rakewire", configFileName)
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
