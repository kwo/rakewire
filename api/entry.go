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

// EntryUpdate updates entries.
func (z *API) EntryUpdate(ctx context.Context, req *msg.EntryUpdateRequest) (*msg.EntryUpdateResponse, error) {

	log.Debugf("here 0")

	user := ctx.Value("user").(*auth.User)
	rsp := &msg.EntryUpdateResponse{}

	err := z.db.Update(func(tx model.Transaction) error {

		log.Debugf("here 1")

		subs := model.S.GetForUser(tx, user.ID)
		subsByFeedID := subs.ByFeedID()
		feedsByURL := model.F.GetBySubscriptions(tx, subs).ByURL()

		log.Debugf("here 2")
		for url, entries := range req.Entries.BySubscription() {
			log.Debugf("here 3")
			if feed, ok := feedsByURL[url]; ok {
				if _, ok := subsByFeedID[feed.ID]; ok { // ignore possibly malicious urls, for which no user subscription
					for _, entry := range entries {
						if item := model.I.GetByGUID(tx, feed.ID, entry.GUID); item != nil {
							if e := model.E.Get(tx, feed.ID, item.ID); e != nil {
								e.Read = entry.Read
								e.Star = entry.Star
								if err := model.E.Save(tx, e); err != nil {
									return err
								} // save
							} // get entry
						} // get item
					} // loop entries
				} // verify subscription
			} // get feed by url
		} // loop url, request entries

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
