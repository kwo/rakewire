package main

import (
	"os"
	"os/signal"
	"path"
	"rakewire.com/db/bolt"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"rakewire.com/server"
	"syscall"
)

const (
	configFileName = "config.yaml"
)

var (
	httpd    *server.Httpd
	database *bolt.Database
	logger   = logging.New("main")
)

func main() {

	cfg := getConfig()
	if cfg == nil {
		logger.Printf("Abort! No config file found at %s\n", getConfigFileLocation())
		os.Exit(1)
		return
	}

	database = &bolt.Database{}
	err := database.Open(&cfg.Database)
	if err != nil {
		logger.Printf("Abort! Cannot access database: %s\n", err.Error())
		os.Exit(1)
		return
	}

	httpd = &server.Httpd{}
	go httpd.Start(cfg.Httpd)

	waitForSignals()

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

func waitForSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logging.Linefeed()
	logger.Println("stopping... ")

	// shutdown server
	_ = httpd.Stop()
	httpd = nil

	// close database
	_ = database.Close()
	database = nil

	logger.Println("done")

}
