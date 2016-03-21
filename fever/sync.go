package fever

import (
	"log"
	"rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getSavedItemIDs(userID string, tx model.Transaction) (string, error) {

	entries := model.E.Query(tx, userID).Starred()
	log.Printf("%-7s %-7s saved count %d", logDebug, logName, len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ItemID), 10))
	}

	return strings.Join(idArray, ","), nil

}

func (z *API) getUnreadItemIDs(userID string, tx model.Transaction) (string, error) {

	entries := model.E.Query(tx, userID).Unread()
	log.Printf("%-7s %-7s unread count %d", logDebug, logName, len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ItemID), 10))
	}

	return strings.Join(idArray, ","), nil

}
