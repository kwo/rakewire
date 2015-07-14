package main

import (
	"os"
	"os/signal"
	"path"
	"rakewire.com/db/bolt"
	"rakewire.com/httpd"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"syscall"
)

const (
	configFileName = "config.yaml"
)

var (
	ws     *httpd.Httpd
	db     *bolt.Database
	logger = logging.New("main")
)

func main() {

	cfg := getConfig()
	if cfg == nil {
		logger.Printf("Abort! No config file found at %s\n", getConfigFileLocation())
		os.Exit(1)
		return
	}

	db = &bolt.Database{}
	err := db.Open(&cfg.Database)
	if err != nil {
		logger.Printf("Abort! Cannot access database: %s\n", err.Error())
		os.Exit(1)
		return
	}

	chErrors := make(chan error)

	ws = &httpd.Httpd{
		Database: db,
	}
	go ws.Start(&cfg.Httpd, chErrors)

	waitForSignals(chErrors)

}

func waitForSignals(chErrors chan error) {
	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-chErrors:
		logger.Printf("Received error: %s", err.Error())
	case <-chSignals:
	}

	logging.Linefeed()
	logger.Println("Stopping... ")

	// shutdown server
	ws.Stop()
	ws = nil

	// close database
	db.Close()
	db = nil

	logger.Println("Done")

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
