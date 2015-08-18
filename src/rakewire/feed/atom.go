package feed

import (
	"github.com/kwo/ocd/feeds/atom"
	"rakewire/logging"
)

var (
	logger = logging.New("feed")
)

func atomToFeed(a *atom.Feed) (*Feed, error) {

	f := &Feed{}

	f.Flavor = "atom"
	f.Generator = atomGeneratorToString(a.Metadata.Generator)
	f.Icon = a.Metadata.Icon
	f.ID = a.Metadata.Id
	f.Links = make(map[string]string)
	f.Rights = a.Metadata.Rights.Text
	f.Subtitle = a.Metadata.Subtitle.Text
	f.Title = a.Title
	f.Updated = &a.Metadata.Updated

	for j := range a.Authors {
		b := &Person{a.Authors[j].Name, a.Authors[j].Email, a.Authors[j].Uri}
		f.Authors = append(f.Authors, b)
	}

	for j := range a.Links {
		f.Links[a.Links[j].Rel] = a.Links[j].Href
	}

	for i := range a.Entries {

		entry := &Entry{}
		f.Entries = append(f.Entries, entry)

		entry.Content = a.Entries[i].Content.Text
		entry.Created = &a.Entries[i].Published
		entry.ID = a.Entries[i].Id
		entry.Links = make(map[string]string)
		entry.Summary = a.Entries[i].Summary.Text
		entry.Title = a.Entries[i].Title
		entry.Updated = &a.Entries[i].Updated

		for j := range a.Entries[i].Authors {
			b := &Person{a.Entries[i].Authors[j].Name, a.Entries[i].Authors[j].Email, a.Entries[i].Authors[j].Uri}
			entry.Authors = append(entry.Authors, b)
		}

		for j := range a.Entries[i].Categories {
			if a.Entries[i].Categories[j].Term != "" {
				entry.Categories = append(entry.Categories, a.Entries[i].Categories[j].Term)
			}
		}

		for j := range a.Entries[i].Links {
			entry.Links[a.Entries[i].Links[j].Rel] = a.Entries[i].Links[j].Href
		}

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

func atomGeneratorToString(g atom.Generator) string {

	var result string

	if g.Text != "" {
		result = g.Text
		if g.Version != "" {
			result += " " + g.Version
		}
		if g.Uri != "" {
			result += " (" + g.Uri + ")"
		}
	}

	return result

}
