package pollfeed

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	"rakewire.com/db/bolt"
	"testing"
	"time"
)

const (
	databaseFile = "../test/pollfeed.db"
)

func TestTickerKillSignal(t *testing.T) {

	beenThere := false
	killsignal := make(chan bool)
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
	run:
		for {
			select {
			case <-ticker.C:
				assert.Fail(t, "ticker should never fire")
			case <-killsignal:
				ticker.Stop()
				break run
			}
		}
		beenThere = true
	}()
	killsignal <- true
	assert.True(t, beenThere)

}

func TestTickerPositive(t *testing.T) {

	beenThere := false
	ticker := time.NewTicker(1 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				beenThere = !beenThere
				ticker.Stop()
				break
			}
		}
	}()
	time.Sleep(2 * time.Millisecond)
	assert.True(t, beenThere)

}

func TestPoll(t *testing.T) {

	//t.SkipNow()

	// open database
	database := &bolt.Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// create service
	cfg := &Configuration{
		FrequencyMinutes: 1,
	}
	pf := NewService(cfg, database)
	pf.pollFrequency = 50 * time.Millisecond

	pf.Start()
	require.Equal(t, true, pf.IsRunning())
	time.Sleep(100 * time.Millisecond)
	pf.Stop()
	assert.Equal(t, false, pf.IsRunning())

	// close database
	err = database.Close()
	assert.Nil(t, err)

	// remove file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
