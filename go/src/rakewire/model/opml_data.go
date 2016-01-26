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

	subscriptions, err := SubscriptionsByUser(user.ID, tx)
	if err != nil {
		return nil, err
	}

	groupsByID := groups.GroupByID()
	subscriptionsByGroup := groupSubscriptionsByGroup(subscriptions, groupsByID)

	categories := make(map[string]*Outline)
	for group, subscriptions1 := range subscriptionsByGroup {

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

		for _, subscription := range subscriptions1 {

			flags := ""
			if subscription.AutoRead {
				flags += " +autoread"
			}
			if subscription.AutoStar {
				flags += " +autostar"
			}
			flags = strings.TrimSpace(flags)

			var created *time.Time
			if !subscription.DateAdded.IsZero() {
				x := subscription.DateAdded.UTC()
				created = &x
			}

			getTitle := func(uf *Subscription) string {
				result := subscription.Title
				if result == "" {
					result = subscription.Feed.Title
				}
				if result == "" {
					result = subscription.Feed.URL
				}
				return result
			}

			outline := &Outline{
				Type:        "rss",
				Text:        getTitle(subscription),
				Title:       getTitle(subscription),
				Created:     created,
				Description: subscription.Notes,
				Category:    flags,
				XMLURL:      subscription.Feed.URL,
				HTMLURL:     subscription.Feed.SiteURL,
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
	groupsByName := groups.GroupByName()
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

	subscriptions, err := SubscriptionsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, subscription := range subscriptions {
		subscription.GroupIDs = []uint64{}
		subscription.AutoRead = false
		subscription.AutoStar = false
	}
	subscriptionsByURL, _ := groupSubscriptionsByURL(subscriptions)

	for branch, outlines := range flatOPML {

		group := groupsByName[branch.Text]

		for _, outline := range outlines {

			uf := subscriptionsByURL[outline.XMLURL]
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

				uf = NewSubscription(userID, f.ID)
				uf.Feed = f
				log.Printf("%-7s %-7s adding subscription: %s", logDebug, logName, uf.Feed.URL)

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

		// remove unused subscriptions
		subscriptions, err := SubscriptionsByUser(userID, tx)
		if err != nil {
			return err
		}
		_, subscriptionDuplicates := groupSubscriptionsByURL(subscriptions)
		for _, subscription := range subscriptions {
			if _, ok := outlinesByURL[subscription.Feed.URL]; !ok {
				log.Printf("%-7s %-7s removing subscription: %s", logDebug, logName, subscription.Feed.URL)
				if err := subscription.Delete(tx); err != nil {
					return err
				}
			}
		}
		for _, subscription := range subscriptionDuplicates {
			log.Printf("%-7s %-7s removing duplicate subscription: %s", logDebug, logName, subscription.Feed.URL)
			if err := subscription.Delete(tx); err != nil {
				return err
			}
		}

		// remove unused groups // TODO: remove unused groups can be done in maintenance thread
		uniqueGroups := collectGroups(subscriptions)
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

func groupSubscriptionsByGroup(subscriptions Subscriptions, groups map[uint64]*Group) map[*Group]Subscriptions {

	result := make(map[*Group]Subscriptions)
	for _, subscription := range subscriptions {
		for _, groupID := range subscription.GroupIDs {
			result[groups[groupID]] = append(result[groups[groupID]], subscription)
		}
	}
	return result

}

func groupSubscriptionsByURL(subscriptions Subscriptions) (map[string]*Subscription, Subscriptions) {
	result := make(map[string]*Subscription)
	duplicates := Subscriptions{}
	for _, subscription := range subscriptions {
		if _, ok := result[subscription.Feed.URL]; !ok {
			result[subscription.Feed.URL] = subscription
		} else {
			duplicates = append(duplicates, subscription)
		}
	}
	return result, duplicates
}

func collectGroups(subscriptions Subscriptions) map[uint64]int {
	result := make(map[uint64]int)
	for _, subscription := range subscriptions {
		for _, groupID := range subscription.GroupIDs {
			result[groupID] = result[groupID] + 1
		}
	}
	return result
}
