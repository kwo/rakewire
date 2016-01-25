package fever

import (
	"log"
	"rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getSavedItemIDs(userID uint64, tx model.Transaction) (string, error) {

	userentries, err := model.UserEntriesStarredByUser(userID, tx)
	if err != nil {
		return "", err
	}

	idArray := []string{}
	for _, userentry := range userentries {
		id := strconv.FormatUint(userentry.ID, 10)
		idArray = append(idArray, id)
	}

	return strings.Join(idArray, ","), nil

}

func (z *API) getUnreadItemIDs(userID uint64, tx model.Transaction) (string, error) {

	userentries, err := model.UserEntriesUnreadByUser(userID, tx)
	if err != nil {
		return "", err
	}
	log.Printf("%-7s %-7s userentry count %d", logDebug, logName, len(userentries))

	idArray := []string{}
	for _, userentry := range userentries {
		id := strconv.FormatUint(userentry.ID, 10)
		idArray = append(idArray, id)
	}

	return strings.Join(idArray, ","), nil

}
