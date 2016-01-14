package opml

import (
	"fmt"
	"log"
	"rakewire/db"
	"rakewire/model"
	"strings"
	"time"
)

const (
	logName  = "[opml]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// TODO: sort entries on export

// Export OPML document
func Export(user *model.User, database db.Database) (*OPML, error) {

	groups, err := database.GroupGetAllByUser(user.ID)
	if err != nil {
		return nil, err
	}

	userfeeds, err := database.UserFeedGetAllByUser(user.ID)
	if err != nil {
		return nil, err
	}

	groupsByID := groupGroupsByID(groups)
	userfeedsByGroup := groupUserFeedsByGroup(userfeeds, groupsByID)

	categories := make(map[string]*Outline)
	for group, userfeeds1 := range userfeedsByGroup {

		log.Printf("%-7s %-7s category: %s", logDebug, logName, group.Name)

		if _, ok := categories[group.Name]; !ok {
			category := &Outline{
				Text:  group.Name,
				Title: group.Name,
			}
			categories[group.Name] = category

		}

		category := categories[group.Name]

		// TODO: if all outlines in a group are autoread/autostar assign to group and remove from outlines

		for _, userfeed := range userfeeds1 {

			flags := ""
			if userfeed.AutoRead {
				flags += " +autoread"
			}
			if userfeed.AutoStar {
				flags += " +autostar"
			}
			flags = strings.TrimSpace(flags)

			var created *time.Time
			if !userfeed.DateAdded.IsZero() {
				x := userfeed.DateAdded.UTC()
				created = &x
			}

			outline := &Outline{
				Type:        "rss",
				Text:        userfeed.Title,
				Title:       userfeed.Title,
				Created:     created,
				Description: userfeed.Notes,
				Category:    flags,
				XMLURL:      userfeed.Feed.URL,
				HTMLURL:     userfeed.Feed.SiteURL,
			}
			log.Printf("%-7s %-7s feed: %s", logDebug, logName, outline.Text)
			category.Outlines = append(category.Outlines, outline)
		}
	}

	outlines := []*Outline{}
	for _, category := range categories {
		log.Printf("%-7s %-7s outline %s: %d", logDebug, logName, category.Text, len(category.Outlines))
		outlines = append(outlines, category)
	}

	opml := &OPML{
		Head: &Head{
			Title:       "Rakewire Subscriptions",
			DateCreated: time.Now().UTC().Truncate(time.Second),
			OwnerName:   user.Username,
		},
		Body: &Body{
			Outlines: outlines,
		},
	}

	return opml, nil

}

// Import OPML document into database
func Import(userID uint64, opml *OPML, replace bool, database db.Database) error {

	flatOPML := Flatten(opml.Body)

	// add missing groups
	groups, err := database.GroupGetAllByUser(userID)
	if err != nil {
		return err
	}
	groupsByName := groupGroupsByName(groups)

	for branch := range flatOPML {
		group := groupsByName[branch.Name]
		if group == nil {
			group = model.NewGroup(userID, branch.Name)
			groupsByName[branch.Name] = group
		}
		if err := database.GroupSave(group); err != nil {
			return err
		}
	}

	userfeeds, err := database.UserFeedGetAllByUser(userID)
	if err != nil {
		return err
	}
	for _, userfeed := range userfeeds {
		userfeed.GroupIDs = []uint64{}
		userfeed.AutoRead = false
		userfeed.AutoStar = false
	}
	userfeedsByURL, _ := groupUserFeedsByURL(userfeeds)

	for branch, outlines := range flatOPML {

		group := groupsByName[branch.Name]
		if group == nil {
			return fmt.Errorf("Group not found: %s", branch.Name)
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
				uf.Feed = f
				log.Printf("%-7s %-7s adding userfeed: %s", logDebug, logName, uf.Feed.URL)

			}

			uf.Title = outline.Title
			uf.Notes = outline.Description
			uf.AutoRead = uf.AutoRead || branch.AutoRead || strings.Contains(outline.Category, "+autoread")
			uf.AutoStar = uf.AutoStar || branch.AutoStar || strings.Contains(outline.Category, "+autostar")
			uf.AddGroup(group.ID)
			if outline.Created != nil && !outline.Created.IsZero() {
				uf.DateAdded = *outline.Created
			}
			if uf.DateAdded.IsZero() {
				uf.DateAdded = time.Now().Truncate(time.Second)
			}
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
		for _, group := range groupsByName {
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

func groupGroupsByID(groups []*model.Group) map[uint64]*model.Group {
	result := make(map[uint64]*model.Group)
	for _, group := range groups {
		result[group.ID] = group
	}
	return result
}

func groupGroupsByName(groups []*model.Group) map[string]*model.Group {
	result := make(map[string]*model.Group)
	for _, group := range groups {
		result[group.Name] = group
	}
	return result
}

func groupUserFeedsByGroup(userfeeds []*model.UserFeed, groups map[uint64]*model.Group) map[*model.Group][]*model.UserFeed {

	result := make(map[*model.Group][]*model.UserFeed)
	for _, userfeed := range userfeeds {
		for _, groupID := range userfeed.GroupIDs {
			result[groups[groupID]] = append(result[groups[groupID]], userfeed)
		}
	}
	return result

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

func groupOutlinesByURL(flatOPML map[*Branch][]*Outline) map[string]*Outline {
	result := make(map[string]*Outline)
	for _, outlines := range flatOPML {
		for _, outline := range outlines {
			if _, ok := result[outline.XMLURL]; !ok {
				result[outline.XMLURL] = outline
			}
		}
	}
	return result
}
