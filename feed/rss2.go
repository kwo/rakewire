package feed

import (
	"encoding/xml"
	"time"
)

type rssFeed struct {
	XMLName xml.Name
	Channel rssChannel `xml:"channel"`
	Version string     `xml:"version,attr"`
}

type rssChannel struct {
	Description string    `xml:"description"`
	Generator   string    `xml:"generator"`
	Items       []rssItem `xml:"item"`
	Link        string    `xml:"rss link"`
	Title       string    `xml:"title"`
}

type rssItem struct {
	Content string    `xml:"content"`
	Created time.Time `xml:"pubDate"`
	ID      string    `xml:"guid"`
	Summary string    `xml:"description"`
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
}

func (r rssFeed) toFeed() (*Feed, error) {

	f := &Feed{}

	f.ID = r.Channel.Link
	f.Title = r.Channel.Title
	f.Subtitle = r.Channel.Description
	f.Flavor = "rss" + r.Version
	f.Generator = r.Channel.Generator

	for _, rssItem := range r.Channel.Items {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.Title = rssItem.Title
		*entry.Created = rssItem.Created
		*entry.Updated = rssItem.Updated
		if entry.Updated.IsZero() {
			entry.Updated = entry.Created
		}
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

		entry.Summary = rssItem.Summary
		entry.Content = rssItem.Content

	} // loop

	if len(f.Entries) > 0 {
		f.Updated = f.Entries[0].Updated
	}

	logger.Printf("RSS SPACE:!%s!", r.XMLName.Space)

	return f, nil

}
