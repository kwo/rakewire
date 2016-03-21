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

		entry := model.E.Get(tx, userID, encodeID(idStr))
		if entry == nil {
			return fmt.Errorf("Entry not found: %s", idStr)
		}

		switch pAs {
		case "read":
			entry.Read = true
		case "unread":
			entry.Read = false
		case "saved":
			entry.Star = true
		case "unsaved":
			entry.Star = false
		default:
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}

		if err := model.E.Save(tx, entry); err != nil {
			return err
		}

	case "feed":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		feedID := encodeID(idStr)
		entries := model.E.Query(tx, userID).Feed(feedID).Max(maxTime).Unread()
		for _, entry := range entries {
			entry.Read = false
		}
		if err := model.E.SaveAll(tx, entries); err != nil {
			return err
		}

	case "group":
		if pAs != "read" {
			return fmt.Errorf("Invalid value for as parameter: %s", pAs)
		}
		groupID := encodeID(idStr)
		subscriptions := model.S.GetForUser(tx, userID).WithGroup(groupID)
		for _, subscription := range subscriptions {
			entries := model.E.Query(tx, userID).Feed(subscription.FeedID).Max(maxTime).Unread()
			for _, entry := range entries {
				entry.Read = false
			}
			if err := model.E.SaveAll(tx, entries); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("Invalid value for mark parameter: %s", mark)
	}

	return nil
}
