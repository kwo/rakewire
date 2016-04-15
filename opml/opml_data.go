package opml

import (
	"fmt"
	"rakewire/model"
	"strings"
	"time"
)

// Export OPML document
func Export(tx model.Transaction, user *model.User) (*OPML, error) {

	groups := model.G.GetForUser(tx, user.ID)
	groupsByID := groups.ByID()
	subscriptions := model.S.GetForUser(tx, user.ID)
	subscriptionsByGroup := groupSubscriptionsByGroup(subscriptions, groupsByID)
	feedsByID := model.F.GetBySubscriptions(tx, subscriptions).ByID()

	categories := make(map[string]*Outline)
	for group, groupSubscriptions := range subscriptionsByGroup {

		if _, ok := categories[group.Name]; !ok {
			category := &Outline{
				Title: group.Name,
			}
			categories[group.Name] = category

		}

		category := categories[group.Name]

		for _, subscription := range groupSubscriptions {

			feed := feedsByID[subscription.FeedID]
			if feed == nil {
				return nil, fmt.Errorf("Missing feed for subscription, feedID: %s", subscription.FeedID)
			}

			flags := ""
			if subscription.AutoRead {
				flags += " +autoread"
			}
			if subscription.AutoStar {
				flags += " +autostar"
			}
			flags = strings.TrimSpace(flags)

			var created *time.Time
			if !subscription.Added.IsZero() {
				x := subscription.Added.UTC()
				created = &x
			}

			getTitle := func(s *model.Subscription) string {
				result := subscription.Title
				if result == "" {
					result = feed.Title
				}
				if result == "" {
					result = feed.URL
				}
				return result
			}

			outline := &Outline{
				Type:        "rss",
				Title:       getTitle(subscription),
				Created:     created,
				Description: subscription.Notes,
				Category:    flags,
				XMLURL:      feed.URL,
				HTMLURL:     feed.SiteURL,
			}
			category.Outlines = append(category.Outlines, outline)
		}
	}

	outlines := Outlines{}
	for _, category := range categories {
		outlines = append(outlines, category)
	}

	outlines.Sort()

	now := time.Now().UTC().Truncate(time.Second)
	opml := &OPML{
		Head: &Head{
			Title:       "Rakewire Subscriptions",
			DateCreated: &now,
			OwnerName:   user.Username,
		},
		Body: &Body{
			Outlines: outlines,
		},
	}

	return opml, nil

}

// Import OPML document into database
func Import(tx model.Transaction, userID string, opml *OPML) error {

	flatOPML := flatten(opml.Body.Outlines)

	// add missing groups
	groups := model.G.GetForUser(tx, userID)
	groupsByName := groups.ByName()
	for branch := range flatOPML {
		group := groupsByName[branch.Title]
		if group == nil {
			group = model.G.New(userID, branch.Title)
			groupsByName[branch.Title] = group
		}
		if err := model.G.Save(tx, group); err != nil {
			return err
		}
	}

	// get subscriptions, reset
	subscriptions := model.S.GetForUser(tx, userID)
	feedsByID := model.F.GetBySubscriptions(tx, subscriptions).ByID()
	for _, subscription := range subscriptions {
		subscription.GroupIDs = []string{}
		subscription.AutoRead = false
		subscription.AutoStar = false
	}
	subscriptionsByURL, _ := groupSubscriptionsByURL(subscriptions, feedsByID)

	for branch, outlines := range flatOPML {

		for _, outline := range outlines {

			var feed *model.Feed
			subscription := subscriptionsByURL[outline.XMLURL]
			if subscription == nil {

				feed = model.F.GetByURL(tx, outline.XMLURL)
				if feed == nil {
					feed = model.F.New(outline.XMLURL)
					if err := model.F.Save(tx, feed); err != nil {
						return err
					}
				}

				subscription = model.S.New(userID, feed.ID)

				subscriptions = append(subscriptions, subscription)
				feedsByID[feed.ID] = feed
				subscriptionsByURL[feed.URL] = subscription

			} // subscription nil

			getTitle := func() string {
				result := outline.Title
				if result == "" {
					result = outline.Text
				}
				if result == "" {
					result = feed.Title
				}
				if result == "" {
					result = feed.SiteURL
				}
				if result == "" {
					result = feed.URL
				}
				return result
			}

			// get group, add if necessary
			group := groupsByName[branch.Title]
			if group == nil {
				group = model.G.New(userID, branch.Title)
				if err := model.G.Save(tx, group); err != nil {
					return err
				}
				groups = append(groups, group)
				groupsByName[group.Name] = group
			}

			subscription.Title = getTitle()
			subscription.Notes = outline.Description
			subscription.AutoRead = subscription.AutoRead || branch.IsAutoRead() || outline.IsAutoRead()
			subscription.AutoStar = subscription.AutoStar || branch.IsAutoStar() || outline.IsAutoStar()
			subscription.AddGroup(group.ID)

			// ignore outline.Created on import, only useful for informational purposes on export
			if subscription.Added.IsZero() {
				subscription.Added = time.Now().Truncate(time.Second)
			}

			if err := model.S.Save(tx, subscription); err != nil {
				return err
			}

		} // outlines

	} // flatOPML

	// remove unused subscriptions
	outlinesByURL := groupOutlinesByURL(flatOPML)
	for url, subscription := range subscriptionsByURL {
		if _, ok := outlinesByURL[url]; !ok {
			if err := model.S.Delete(tx, subscription.GetID()); err != nil {
				return err
			}
		}
	}

	return nil

}
