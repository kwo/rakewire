package fever

import (
	"log"
	"strconv"
	"strings"
)

func (z *API) getSavedItemIDs(userID uint64) (string, error) {

	userentries, err := z.db.UserEntryGetStarredForUser(userID)
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

func (z *API) getUnreadItemIDs(userID uint64) (string, error) {

	userentries, err := z.db.UserEntryGetNext(userID, 0, 0)
	if err != nil {
		return "", err
	}
	log.Printf("%-7s %-7s userentry count %d", logDebug, logName, len(userentries))
	for i, ue := range userentries {
		log.Printf("%-7s %-7s userentry %d: %v", logDebug, logName, i, ue)
	}

	userentries, err = z.db.UserEntryGetUnreadForUser(userID)
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
