package feed

import (
	"time"
)

type rssFeed struct {
	Channel rssChannel `xml:"channel"`
	Version string     `xml:"version,attr"`
}

type rssChannel struct {
	Description string `xml:"description"`
	//Items       []rssItem `xml:"item"`
	Link  string `xml:"link"`
	Title string `xml:"title"`
}

type rssItem struct {
	ID         string         `xml:"id"`
	Published  time.Time      `xml:"published"`
	Updated    time.Time      `xml:"updated"`
	Author     atomAuthor     `xml:"author"`
	Title      string         `xml:"title"`
	Categories []atomCategory `xml:"category"`
	Links      []atomLink     `xml:"link"`
	Summary    atomText       `xml:"summary"`
	Content    atomText       `xml:"content"`
}

func (r rssChannel) toFeed() (*Feed, error) {

	var f Feed

	f.ID = r.Link
	f.Title = r.Title
	f.Subtitle = r.Description
	// f.Author = &Author{r.Author.Name, r.Author.EMail, r.Author.URI}
	// f.Icon = r.Icon
	// f.Rights = r.Rights
	// f.Generator = r.Generator.String()

	// if !a.Updated.IsZero() {
	// 	f.Updated = &a.Updated
	// }
	//
	// for j := 0; j < len(a.Links); j++ {
	// 	atomLink := a.Links[j]
	// 	link := Link{atomLink.Rel, atomLink.Href}
	// 	f.Links = append(f.Links, &link)
	// }
	//
	// for i := 0; i < len(a.Entries); i++ {
	//
	// 	entry := Entry{}
	// 	atomEntry := a.Entries[i]
	// 	f.Entries = append(f.Entries, &entry)
	//
	// 	entry.ID = atomEntry.ID
	// 	entry.Title = atomEntry.Title
	// 	if !atomEntry.Created.IsZero() {
	// 		entry.Created = &atomEntry.Created
	// 	}
	// 	if !atomEntry.Updated.IsZero() {
	// 		entry.Updated = &atomEntry.Updated
	// 	}
	// 	entry.Author = &Author{atomEntry.Author.Name, atomEntry.Author.EMail, atomEntry.Author.URI}
	//
	// 	for j := 0; j < len(atomEntry.Links); j++ {
	// 		atomLink := atomEntry.Links[j]
	// 		link := Link{atomLink.Rel, atomLink.Href}
	// 		entry.Links = append(entry.Links, &link)
	// 	}
	//
	// 	for j := 0; j < len(atomEntry.Categories); j++ {
	// 		entry.Categories = append(entry.Categories, atomEntry.Categories[j].String())
	// 	}
	//
	// 	entry.Summary = atomEntry.Summary.String()
	// 	entry.Content = atomEntry.Content.String()
	//
	// }

	return &f, nil

}
