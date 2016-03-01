package fever

import (
	"log"
	"rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getSavedItemIDs(userID string, tx model.Transaction) (string, error) {

	entries, err := model.EntriesStarredByUser(userID, tx)
	if err != nil {
		return "", err
	}
	log.Printf("%-7s %-7s saved count %d", logDebug, logName, len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ID), 10))
	}

	return strings.Join(idArray, ","), nil

}

func (z *API) getUnreadItemIDs(userID string, tx model.Transaction) (string, error) {

	entries, err := model.EntriesUnreadByUser(userID, tx)
	if err != nil {
		return "", err
	}
	log.Printf("%-7s %-7s unread count %d", logDebug, logName, len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ID), 10))
	}

	return strings.Join(idArray, ","), nil

}
