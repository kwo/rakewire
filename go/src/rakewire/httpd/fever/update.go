package fever

import (
	"fmt"
	"rakewire/model"
	"strconv"
	"time"
)

func (z *API) updateItems(userID uint64, mark, pAs, idStr, beforeStr string, tx model.Transaction) error {

	if idStr == "-1" {
		return fmt.Errorf("sparks not supported")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
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

		items, err := model.UserEntriesByUser(userID, []uint64{id}, tx)
		if err != nil {
			return err
		}
		if len(items) != 1 {
			return fmt.Errorf("User entry not found: %s", idStr)
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

		if err := model.UserEntriesSave([]*model.UserEntry{item}, tx); err != nil {
			return err
		}

	case "feed":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		if err := model.UserEntriesUpdateReadByFeed(userID, id, maxTime, true, tx); err != nil {
			return err
		}

	case "group":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		if err := model.UserEntriesUpdateReadByGroup(userID, id, maxTime, true, tx); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Invalid value for mark parameter: %s", mark)
	}

	return nil
}
