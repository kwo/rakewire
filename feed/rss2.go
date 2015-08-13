package feed

import (
	"github.com/kwo/ocd/feeds/rss"
	"time"
)

func rssToFeed(r *rss.Rss) (*Feed, error) {

	f := &Feed{}

	f.ID = r.Channel.Link
	f.Title = r.Channel.Title
	f.Subtitle = r.Channel.Description
	f.Flavor = "rss2"
	f.Generator = r.Channel.Generator
	f.Updated = getTime(r.Channel.PubDate)

	for i := range r.Channel.Items {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.Title = r.Channel.Items[i].Title
		entry.Updated = getTime(r.Channel.Items[i].PubDate)
		//entry.Created = rssItem.Created
		//entry.Updated = rssItem.Updated
		// if entry.Updated.IsZero() {
		// 	entry.Updated = entry.Created
		// }
		//entry.Author = &Author{atomEntry.Author.Name, atomEntry.Author.EMail, atomEntry.Author.URI}

		// for j := 0; j < len(atomEntry.Links); j++ {
		// 	atomLink := atomEntry.Links[j]
		// 	link := Link{atomLink.Rel, atomLink.Href}
		// 	entry.Links = append(entry.Links, &link)
		// }
		//
		// for j := 0; j < len(atomEntry.Categories); j++ {
		// 	entry.Categories = append(entry.Categories, atomEntry.Categories[j].String())
		// }

		// entry.Summary = rssItem.Summary
		// entry.Content = rssItem.Content

	} // loop

	// set updated to first entry time ignoring time in header
	if len(f.Entries) > 0 && f.Entries[0].Updated != nil && !f.Entries[0].Updated.IsZero() {
		f.Updated = f.Entries[0].Updated
	}

	// turn zero times to nil
	if f.Updated != nil && f.Updated.IsZero() {
		f.Updated = nil
	}

	return f, nil

}

func getTime(t rss.Time) *time.Time {
	dt := t.GetTime()
	if dt.IsZero() {
		return nil
	}
	return &dt
}
