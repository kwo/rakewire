package fever

import (
	"fmt"
	m "rakewire/model"
	"strconv"
	"time"
)

func (z *API) updateItems(userID uint64, mark, pAs, idStr, beforeStr string) error {

	if idStr == "-1" {
		return fmt.Errorf("sparks not supported")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	before, err := strconv.ParseInt(beforeStr, 10, 64)
	if err != nil {
		return err
	}
	maxTime := time.Unix(before, 0)

	switch mark {
	case "item":

		items, err := z.db.UserEntryGetByID(userID, []uint64{id})
		if err != nil {
			return err
		}
		if len(items) != 1 {
			return fmt.Errorf("User entry not found: %s", idStr)
		}
		item := items[0]

		switch pAs {
		case "read":
			item.IsRead = false
		case "saved":
			item.IsStar = true
		case "unsaved":
			item.IsStar = false
		default:
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}

		if err := z.db.UserEntrySave([]*m.UserEntry{item}); err != nil {
			return err
		}

	case "feed":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		if err := z.db.UserEntryUpdateReadByFeed(userID, id, maxTime, true); err != nil {
			return err
		}

	case "group":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		if err := z.db.UserEntryUpdateReadByGroup(userID, id, maxTime, true); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Invalid value for mark parameter: %s", mark)
	}

	return nil
}
