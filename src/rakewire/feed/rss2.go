package feed

import (
	"github.com/kwo/ocd/feeds/rss"
	"time"
)

func rssToFeed(r *rss.Rss) (*Feed, error) {

	f := &Feed{}

	f.Flavor = "rss2"
	f.Generator = r.Channel.Generator
	f.Icon = r.Channel.Image.Url
	f.ID = r.Channel.Link
	f.Links = make(map[string]string)
	f.Rights = r.Channel.Copyright
	f.Subtitle = r.Channel.Description
	f.Title = r.Channel.Title
	f.Updated = getTime(r.Channel.PubDate.GetTime())

	for i := range r.Channel.Items {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.Content = r.Channel.Items[i].Description
		if r.Channel.Items[i].Encoded != "" {
			entry.Summary = entry.Content
			entry.Content = r.Channel.Items[i].Encoded
		}

		entry.Created = getTime(r.Channel.Items[i].PubDate.GetTime())
		entry.ID = r.Channel.Items[i].Guid.Text
		entry.Links = make(map[string]string)
		entry.Title = r.Channel.Items[i].Title
		entry.Updated = getTime(r.Channel.Items[i].PubDate.GetTime())
		// use the dublincore date in no pubDate
		if !r.Channel.Items[i].Date.IsZero() {
			entry.Updated = &r.Channel.Items[i].Date
		}

		if r.Channel.Items[i].Author != "" {
			b := &Person{Name: r.Channel.Items[i].Author}
			entry.Authors = append(entry.Authors, b)
		} else if r.Channel.Items[i].Creator != "" { // dublincore
			b := &Person{Name: r.Channel.Items[i].Author}
			entry.Authors = append(entry.Authors, b)
		}

		for j := range r.Channel.Items[i].Categories {
			if r.Channel.Items[i].Categories[j].Text != "" {
				entry.Categories = append(entry.Categories, r.Channel.Items[i].Categories[j].Text)
			}
		}

		if r.Channel.Items[i].Link != "" {
			// #DOING:0 which link is the rss entry link
			// entry.Links[""] = r.Channel.Items[i].Link
		}

	} // loop

	// set updated to first entry time ignoring time in header
	if len(f.Entries) > 0 && !f.Entries[0].Updated.IsZero() {
		f.Updated = f.Entries[0].Updated
	}

	return f, nil

}

func getTime(t time.Time) *time.Time {
	return &t
}
