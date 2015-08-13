package feed

import (
	"github.com/kwo/ocd/feeds/rss"
)

func rssToFeed(r *rss.Rss) (*Feed, error) {

	f := &Feed{}

	f.ID = r.Channel.Link
	f.Title = r.Channel.Title
	f.Subtitle = r.Channel.Description
	f.Flavor = "rss2"
	f.Generator = r.Channel.Generator

	for i := range r.Channel.Items {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.Title = r.Channel.Items[i].Title
		dt := r.Channel.Items[i].PubDate.GetTime()
		entry.Updated = &dt
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

	if len(f.Entries) > 0 {
		f.Updated = f.Entries[0].Updated
	}

	return f, nil

}
