package bolt

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"rakewire/db"
	"rakewire/logging"
	"strings"
	"testing"
)

const (
	feedFile     = "../../../../test/feedlist.txt"
	databaseFile = "../../../../test/bolt.db"
)

func TestMain(m *testing.M) {

	// initialize logging
	logging.Init(&logging.Configuration{
		Level: "debug",
	})

	logger.Debug("Logging configured")

	m.Run()

}

func TestInterface(t *testing.T) {

	var d db.Database = &Database{}
	assert.NotNil(t, d)

}

func readFile(feedfile string) ([]string, error) {

	var result []string

	f, err1 := os.Open(feedfile)
	if err1 != nil {
		return nil, err1
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			result = append(result, url)
		}
	}
	f.Close()

	return result, nil

}
