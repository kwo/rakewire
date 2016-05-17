package api

import (
	"github.com/kwo/rakewire/api/msg"
	"github.com/kwo/rakewire/auth"
	"github.com/kwo/rakewire/model"
	"golang.org/x/net/context"
	"strings"
	"time"
)

// SubscriptionAddUpdate adds or updates a subscription
func (z *API) SubscriptionAddUpdate(ctx context.Context, req *msg.SubscriptionAddUpdateRequest) (*msg.SubscriptionAddUpdateResponse, error) {

	user := ctx.Value("user").(*auth.User)

	rsp := &msg.SubscriptionAddUpdateResponse{}

	err := z.db.Update(func(tx model.Transaction) error {

		feed := model.F.GetByURL(tx, req.Subscription.URL)
		if feed == nil {
			feed = model.F.New(req.Subscription.URL)
			if err := model.F.Save(tx, feed); err != nil {
				return err
			}
		}

		var subscription *model.Subscription
		subscriptions := model.S.GetForUser(tx, user.ID).ByFeedID()[feed.ID]
		if len(subscriptions) > 0 {
			subscription = subscriptions[0]
		}
		if subscription == nil {
			subscription = model.S.New(user.ID, feed.ID)
		}

		subscription.Title = req.Subscription.Title
		if len(subscription.Title) == 0 {
			subscription.Title = feed.Title
		}
		if len(subscription.Title) == 0 {
			subscription.Title = feed.URL
		}
		subscription.Notes = req.Subscription.Notes
		subscription.AutoRead = req.Subscription.AutoRead
		subscription.AutoStar = req.Subscription.AutoStar
		if subscription.Added.IsZero() {
			subscription.Added = time.Now().Truncate(time.Second)
		}

		groupsByName := model.G.GetForUser(tx, user.ID).ByName()
		for _, groupName := range req.Subscription.Groups {

			group := groupsByName[groupName]
			if group == nil && req.AddGroups {
				group = model.G.New(user.ID, groupName)
				if err := model.G.Save(tx, group); err != nil {
					return err
				}
			}

			if group != nil {
				subscription.GroupIDs = append(subscription.GroupIDs, group.GetID())
			}

		}

		if len(subscription.GroupIDs) == 0 {
			rsp.Status = msg.StatusErr
			rsp.Message = "Subscriptions must be assigned to at least one group"
			return errEscape
		}

		return model.S.Save(tx, subscription)

	})

	if err != nil && err != errEscape {
		rsp.Status = msg.StatusErr
		rsp.Message = err.Error()
	}

	return rsp, nil

}

// SubscriptionList lists a user's subscrptions.
func (z *API) SubscriptionList(ctx context.Context, req *msg.SubscriptionListRequest) (*msg.SubscriptionListResponse, error) {

	user := ctx.Value("user").(*auth.User)

	rsp := &msg.SubscriptionListResponse{}

	err := z.db.Select(func(tx model.Transaction) error {

		subs := model.S.GetForUser(tx, user.ID)
		feedsByID := model.F.GetBySubscriptions(tx, subs).ByID()
		groups := model.G.GetForUser(tx, user.ID)

		for _, sub := range subs {
			groupNames := []string{}
			for _, group := range groups.WithIDs(sub.GroupIDs...) {
				groupNames = append(groupNames, group.Name)
			}
			subscription := &msg.Subscription{
				URL:      feedsByID[sub.FeedID].URL,
				Title:    sub.Title,
				Groups:   groupNames,
				Notes:    sub.Notes,
				Added:    sub.Added,
				AutoRead: sub.AutoRead,
				AutoStar: sub.AutoStar,
			}
			if len(req.Filter) == 0 || matchFilter(req.Filter, subscription) {
				rsp.Subscriptions = append(rsp.Subscriptions, subscription)
			}
		}

		return nil

	})

	return rsp, err

}
func matchFilter(filter string, subscription *msg.Subscription) bool {

	match := false
	filter = strings.ToLower(filter)

	if !match && strings.Contains(strings.ToLower(subscription.Title), filter) {
		match = true
	}

	if !match && strings.Contains(strings.ToLower(subscription.URL), filter) {
		match = true
	}

	if !match && strings.Contains(strings.ToLower(subscription.Notes), filter) {
		match = true
	}

	if !match && strings.Contains(subscription.Added.Format(time.RFC3339), filter) {
		match = true
	}

	if !match {
		for _, group := range subscription.Groups {
			if strings.Contains(strings.ToLower(group), filter) {
				match = true
			}
		}
	}

	return match

}

// SubscriptionUnsubscribe removes a user's subscription to a feed
func (z *API) SubscriptionUnsubscribe(ctx context.Context, req *msg.UnsubscribeRequest) (*msg.UnsubscribeResponse, error) {

	user := ctx.Value("user").(*auth.User)

	rsp := &msg.UnsubscribeResponse{}

	err := z.db.Update(func(tx model.Transaction) error {
		feed := model.F.GetByURL(tx, req.URL)
		if feed != nil {
			subscriptions := model.S.GetForUser(tx, user.ID).ByFeedID()[feed.ID]
			for _, subscription := range subscriptions {
				if err := model.S.Delete(tx, subscription.GetID()); err != nil {
					return err
				}
			}
		} else {
			rsp.Status = msg.StatusNotFound
		}
		return nil
	})

	if err != nil && err != errEscape {
		rsp.Status = msg.StatusErr
		rsp.Message = err.Error()
	}

	return rsp, nil

}
