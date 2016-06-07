package api

import (
	"github.com/kwo/rakewire/api/msg"
	"github.com/kwo/rakewire/auth"
	"github.com/kwo/rakewire/model"
	"golang.org/x/net/context"
)

// EntryList lists a subscription's entries.
func (z *API) EntryList(ctx context.Context, req *msg.EntryListRequest) (*msg.EntryListResponse, error) {

	user := ctx.Value("user").(*auth.User)

	rsp := &msg.EntryListResponse{}

	err := z.db.Select(func(tx model.Transaction) error {

		if feed := model.F.GetByURL(tx, req.Subscription); feed != nil {
			if sub := model.S.GetForUser(tx, user.ID).ByFeedID()[feed.GetID()]; sub != nil {

				entries := model.E.Query(tx, user.ID).Feed(feed.ID).Get()
				itemsByID := model.I.GetByEntries(tx, entries).ByID()

				for _, entry := range entries {
					rsp.Entries = append(rsp.Entries, toEntry(entry, itemsByID[entry.ItemID], feed))
				}

			}
		}

		return nil

	})

	return rsp, err

}

func toEntry(entry *model.Entry, item *model.Item, feed *model.Feed) *msg.Entry {

	e := &msg.Entry{}

	e.Subscription = feed.URL
	e.GUID = item.GUID
	e.Title = item.Title
	e.Updated = item.Updated
	e.Read = entry.Read
	e.Star = entry.Star

	return e

}
