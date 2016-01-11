package opml

import (
	"fmt"
	"log"
	"rakewire/db"
	"rakewire/model"
)

const (
	logName  = "[opml]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// Import OPML document into database
func Import(userID uint64, opml *OPML, replace bool, database db.Database) error {

	flatOPML := Flatten(opml.Body)

	// add missing groups
	groups, err := database.GroupGetAllByUser(userID)
	if err != nil {
		return err
	}
	for groupName := range flatOPML {
		group := groups[groupName]
		if group == nil {
			group = model.NewGroup(userID, groupName)
			if err := database.GroupSave(group); err != nil {
				return err
			}
			groups[groupName] = group
		}
	}

	userfeeds, err := database.UserFeedGetAllByUser(userID)
	if err != nil {
		return err
	}
	for _, userfeed := range userfeeds {
		userfeed.GroupIDs = []uint64{}
	}
	userfeedsByURL, _ := groupUserFeedsByURL(userfeeds)

	for groupName, outlines := range flatOPML {

		group := groups[groupName]
		if group == nil {
			return fmt.Errorf("Group not found: %s", groupName)
		}

		for _, outline := range outlines {

			uf := userfeedsByURL[outline.XMLURL]
			if uf == nil {
				f, err := database.GetFeedByURL(outline.XMLURL)
				if err != nil {
					return err
				}
				if f == nil {
					f = model.NewFeed(outline.XMLURL)
					log.Printf("%-7s %-7s adding feed: %s", logDebug, logName, f.URL)
					if _, err := database.FeedSave(f); err != nil {
						return err
					}
				}
				uf = model.NewUserFeed(userID, f.ID)
				uf.Title = outline.Text
				uf.Feed = f
				log.Printf("%-7s %-7s adding userfeed: %s", logDebug, logName, uf.Feed.URL)
				if err := database.UserFeedSave(uf); err != nil {
					return err
				}
			}

			uf.AddGroup(group.ID)
			if err := database.UserFeedSave(uf); err != nil {
				return err
			}

		} // outlines

	} // flatOPML

	if replace {

		outlinesByURL := groupOutlinesByURL(flatOPML)

		// remove unused userfeeds
		userfeeds, err := database.UserFeedGetAllByUser(userID)
		if err != nil {
			return err
		}
		_, userfeedDuplicates := groupUserFeedsByURL(userfeeds)
		for _, userfeed := range userfeeds {
			if _, ok := outlinesByURL[userfeed.Feed.URL]; !ok {
				log.Printf("%-7s %-7s removing userfeed: %s", logDebug, logName, userfeed.Feed.URL)
				if err := database.UserFeedDelete(userfeed); err != nil {
					return err
				}
			}
		}
		for _, userfeed := range userfeedDuplicates {
			log.Printf("%-7s %-7s removing duplicate userfeed: %s", logDebug, logName, userfeed.Feed.URL)
			if err := database.UserFeedDelete(userfeed); err != nil {
				return err
			}
		}

		// remove unused groups
		uniqueGroups := collectGroups(userfeeds)
		for _, group := range groups {
			if _, ok := uniqueGroups[group.ID]; !ok {
				log.Printf("%-7s %-7s removing group: %s", logDebug, logName, group.Name)
				if err := database.GroupDelete(group); err != nil {
					return err
				}
			}
		}

	}

	return nil

}

func groupUserFeedsByURL(userfeeds []*model.UserFeed) (map[string]*model.UserFeed, []*model.UserFeed) {
	result := make(map[string]*model.UserFeed)
	duplicates := []*model.UserFeed{}
	for _, userfeed := range userfeeds {
		if _, ok := result[userfeed.Feed.URL]; !ok {
			result[userfeed.Feed.URL] = userfeed
		} else {
			duplicates = append(duplicates, userfeed)
		}
	}
	return result, duplicates
}

func groupFeedsByURL(feeds []*model.Feed) (map[string]*model.Feed, []*model.Feed) {
	result := make(map[string]*model.Feed)
	duplicates := []*model.Feed{}
	for _, feed := range feeds {
		if _, ok := result[feed.URL]; !ok {
			result[feed.URL] = feed
		} else {
			duplicates = append(duplicates, feed)
		}
	}
	return result, duplicates
}

func collectGroups(userfeeds []*model.UserFeed) map[uint64]int {
	result := make(map[uint64]int)
	for _, userfeed := range userfeeds {
		for _, groupID := range userfeed.GroupIDs {
			result[groupID] = result[groupID] + 1
		}
	}
	return result
}

func groupOutlinesByURL(flatOPML map[string][]Outline) map[string]Outline {
	result := make(map[string]Outline)
	for _, outlines := range flatOPML {
		for _, outline := range outlines {
			if _, ok := result[outline.XMLURL]; !ok {
				result[outline.XMLURL] = outline
			}
		}
	}
	return result
}
