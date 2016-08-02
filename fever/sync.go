package fever

import (
	"strconv"
	"strings"

	"github.com/kwo/rakewire/model"
)

func (z *API) getSavedItemIDs(userID string, tx model.Transaction) (string, error) {

	entries := model.E.Query(tx, userID).Starred()
	log.Debugf("saved count %d", len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ItemID), 10))
	}

	return strings.Join(idArray, ","), nil

}

func (z *API) getUnreadItemIDs(userID string, tx model.Transaction) (string, error) {

	entries := model.E.Query(tx, userID).Unread()
	log.Debugf("unread count %d", len(entries))

	idArray := []string{}
	for _, entry := range entries {
		idArray = append(idArray, strconv.FormatUint(parseID(entry.ItemID), 10))
	}

	return strings.Join(idArray, ","), nil

}
