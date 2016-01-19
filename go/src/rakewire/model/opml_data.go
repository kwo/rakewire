package model

import (
	"log"
	"strings"
	"time"
)

// OPMLExport OPML document
func OPMLExport(user *User, tx Transaction) (*OPML, error) {

	groups, err := GroupsByUser(user.ID, tx)
	if err != nil {
		return nil, err
	}

	userfeeds, err := UserFeedsByUser(user.ID, tx)
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

			getTitle := func(uf *UserFeed) string {
				result := userfeed.Title
				if result == "" {
					result = userfeed.Feed.Title
				}
				if result == "" {
					result = userfeed.Feed.URL
				}
				return result
			}

			outline := &Outline{
				Type:        "rss",
				Text:        getTitle(userfeed),
				Title:       getTitle(userfeed),
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

	outlines := Outlines{}
	for _, category := range categories {
		log.Printf("%-7s %-7s outline %s: %d", logDebug, logName, category.Text, len(category.Outlines))
		outlines = append(outlines, category)
	}

	outlines.Sort()

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

// OPMLImport OPML document into database
func OPMLImport(userID uint64, opml *OPML, replace bool, tx Transaction) error {

	log.Printf("%-7s %-7s importing opml for user %d, replace: %t", logDebug, logName, userID, replace)

	flatOPML := OPMLFlatten(opml.Body.Outlines)

	// add missing groups
	groups, err := GroupsByUser(userID, tx)
	if err != nil {
		return err
	}
	groupsByName := groupGroupsByName(groups)
	for branch := range flatOPML {
		group := groupsByName[branch.Text]
		if group == nil {
			group = NewGroup(userID, branch.Text)
			groupsByName[branch.Text] = group
		}
		if err := group.Save(tx); err != nil {
			return err
		}
	}

	userfeeds, err := UserFeedsByUser(userID, tx)
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

		group := groupsByName[branch.Text]

		for _, outline := range outlines {

			uf := userfeedsByURL[outline.XMLURL]
			if uf == nil {

				f, err := FeedByURL(outline.XMLURL, tx)
				if err != nil {
					return err
				}
				if f == nil {
					f = NewFeed(outline.XMLURL)
					log.Printf("%-7s %-7s adding feed: %s", logDebug, logName, f.URL)
					if _, err := f.Save(tx); err != nil {
						return err
					}
				}

				uf = NewUserFeed(userID, f.ID)
				uf.Feed = f
				log.Printf("%-7s %-7s adding userfeed: %s", logDebug, logName, uf.Feed.URL)

			}

			getTitle := func() string {
				result := outline.Title
				if result == "" {
					result = outline.Text
				}
				if result == "" {
					result = uf.Feed.Title
				}
				if result == "" {
					result = outline.HTMLURL
				}
				if result == "" {
					result = outline.XMLURL
				}
				return result
			}

			uf.Title = getTitle()
			uf.Notes = outline.Description
			uf.AutoRead = uf.AutoRead || branch.IsAutoRead() || outline.IsAutoRead()
			uf.AutoStar = uf.AutoStar || branch.IsAutoStar() || outline.IsAutoStar()
			uf.AddGroup(group.ID)
			if outline.Created != nil && !outline.Created.IsZero() {
				uf.DateAdded = *outline.Created
			}
			if uf.DateAdded.IsZero() {
				uf.DateAdded = time.Now().Truncate(time.Second)
			}
			if err := uf.Save(tx); err != nil {
				return err
			}

		} // outlines

	} // flatOPML

	if replace {

		outlinesByURL := groupOutlinesByURL(flatOPML)

		// remove unused userfeeds
		userfeeds, err := UserFeedsByUser(userID, tx)
		if err != nil {
			return err
		}
		_, userfeedDuplicates := groupUserFeedsByURL(userfeeds)
		for _, userfeed := range userfeeds {
			if _, ok := outlinesByURL[userfeed.Feed.URL]; !ok {
				log.Printf("%-7s %-7s removing userfeed: %s", logDebug, logName, userfeed.Feed.URL)
				if err := userfeed.Delete(tx); err != nil {
					return err
				}
			}
		}
		for _, userfeed := range userfeedDuplicates {
			log.Printf("%-7s %-7s removing duplicate userfeed: %s", logDebug, logName, userfeed.Feed.URL)
			if err := userfeed.Delete(tx); err != nil {
				return err
			}
		}

		// remove unused groups // TODO: remove unused groups can be done in maintenance thread
		uniqueGroups := collectGroups(userfeeds)
		for _, group := range groupsByName {
			if _, ok := uniqueGroups[group.ID]; !ok {
				log.Printf("%-7s %-7s removing group: %s", logDebug, logName, group.Name)
				if err := group.Delete(tx); err != nil {
					return err
				}
			}
		}

	}

	return nil

}

func groupGroupsByID(groups []*Group) map[uint64]*Group {
	result := make(map[uint64]*Group)
	for _, group := range groups {
		result[group.ID] = group
	}
	return result
}

func groupGroupsByName(groups []*Group) map[string]*Group {
	result := make(map[string]*Group)
	for _, group := range groups {
		result[group.Name] = group
	}
	return result
}

func groupUserFeedsByGroup(userfeeds []*UserFeed, groups map[uint64]*Group) map[*Group][]*UserFeed {

	result := make(map[*Group][]*UserFeed)
	for _, userfeed := range userfeeds {
		for _, groupID := range userfeed.GroupIDs {
			result[groups[groupID]] = append(result[groups[groupID]], userfeed)
		}
	}
	return result

}

func groupUserFeedsByURL(userfeeds []*UserFeed) (map[string]*UserFeed, []*UserFeed) {
	result := make(map[string]*UserFeed)
	duplicates := []*UserFeed{}
	for _, userfeed := range userfeeds {
		if _, ok := result[userfeed.Feed.URL]; !ok {
			result[userfeed.Feed.URL] = userfeed
		} else {
			duplicates = append(duplicates, userfeed)
		}
	}
	return result, duplicates
}

func groupFeedsByURL(feeds []*Feed) (map[string]*Feed, []*Feed) {
	result := make(map[string]*Feed)
	duplicates := []*Feed{}
	for _, feed := range feeds {
		if _, ok := result[feed.URL]; !ok {
			result[feed.URL] = feed
		} else {
			duplicates = append(duplicates, feed)
		}
	}
	return result, duplicates
}

func collectGroups(userfeeds []*UserFeed) map[uint64]int {
	result := make(map[uint64]int)
	for _, userfeed := range userfeeds {
		for _, groupID := range userfeed.GroupIDs {
			result[groupID] = result[groupID] + 1
		}
	}
	return result
}
