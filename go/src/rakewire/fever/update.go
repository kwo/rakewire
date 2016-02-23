package fever

import (
	"fmt"
	"rakewire/model"
	"strconv"
	"time"
)

func (z *API) updateItems(userID string, mark, pAs, idStr, beforeStr string, tx model.Transaction) error {

	if idStr == "-1" {
		return fmt.Errorf("sparks not supported")
	}

	maxTime := time.Time{}
	if beforeStr != "" {
		before, err := strconv.ParseInt(beforeStr, 10, 64)
		if err != nil {
			return err
		}
		maxTime = time.Unix(before, 0)
	}

	switch mark {
	case "item":

		items, err := model.EntriesByUser(userID, []string{encodeID(idStr)}, tx)
		if err != nil {
			return err
		}
		if len(items) != 1 {
			return fmt.Errorf("User item not found: %s", idStr)
		}
		item := items[0]

		switch pAs {
		case "read":
			item.IsRead = true
		case "unread":
			item.IsRead = false
		case "saved":
			item.IsStar = true
		case "unsaved":
			item.IsStar = false
		default:
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}

		if err := model.EntriesSave([]*model.Entry{item}, tx); err != nil {
			return err
		}

	case "feed":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		// TODO: first query then use EntriesSave to mark
		if err := model.EntriesUpdateReadByFeed(userID, encodeID(idStr), maxTime, true, tx); err != nil {
			return err
		}

	case "group":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		// TODO: first query then use EntriesSave to mark
		if err := model.EntriesUpdateReadByGroup(userID, encodeID(idStr), maxTime, true, tx); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Invalid value for mark parameter: %s", mark)
	}

	return nil
}
