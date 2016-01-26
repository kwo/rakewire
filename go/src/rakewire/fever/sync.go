package fever

import (
	"log"
	"rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getSavedItemIDs(userID uint64, tx model.Transaction) (string, error) {

	entries, err := model.EntriesStarredByUser(userID, tx)
	if err != nil {
		return "", err
	}

	idArray := []string{}
	for _, entry := range entries {
		id := strconv.FormatUint(entry.ID, 10)
		idArray = append(idArray, id)
	}

	return strings.Join(idArray, ","), nil

}

func (z *API) getUnreadItemIDs(userID uint64, tx model.Transaction) (string, error) {

	entries, err := model.EntriesUnreadByUser(userID, tx)
	if err != nil {
		return "", err
	}
	log.Printf("%-7s %-7s entry count %d", logDebug, logName, len(entries))

	idArray := []string{}
	for _, entry := range entries {
		id := strconv.FormatUint(entry.ID, 10)
		idArray = append(idArray, id)
	}

	return strings.Join(idArray, ","), nil

}
