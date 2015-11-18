package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/db"
	//m "rakewire/model"
	"os"
	"testing"
)

func TestSerialization(t *testing.T) {

	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
