package rss

import (
	"encoding/xml"
	"github.com/kwo/ocd/feeds/modules/atom"
	"github.com/kwo/ocd/feeds/modules/content"
	"github.com/kwo/ocd/feeds/modules/dublincore"
	"github.com/kwo/ocd/feeds/modules/media"
	"strings"
	"time"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string     `xml:"title"`
	Link          string     `xml:"rss link"`
	Description   string     `xml:"description"`
	Language      string     `xml:"language"`
	Copyright     string     `xml:"copyright"`
	WebMaster     string     `xml:"webMaster"`
	PubDate       Time       `xml:"pubDate"`
	LastBuildDate Time       `xml:"lastBuildDate"`
	Categories    []Category `xml:"category"`
	Generator     string     `xml:"generator"`
	Docs          string     `xml:"docs"`
	Cloud         Cloud      `xml:"cloud"`
	Ttl           string     `xml:"ttl"`
	Image         Image      `xml:"image"`
	Rating        string     `xml:"rating"`
	TextInput     TextInput  `xml:"textInput"`
	SkipHours     []int      `xml:"skipHours>hour"`
	SkipDays      []string   `xml:"skipDays>day"`
	Items         []Item     `xml:"item"`
	atom.Atom
}

type Category struct {
	Domain string `xml:"domain,attr"`
	Text   string `xml:",chardata"`
}

type Cloud struct {
	Domain            string `xml:"domain,attr"`
	Port              string `xml:"port,attr"`
	Path              string `xml:"path,attr"`
	RegisterProcedure string `xml:"registerProcedure,attr"`
	Protocol          string `xml:"protocol,attr"`
}

type Image struct {
	Url         string `xml:"url"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Width       string `xml:"width"`
	Height      string `xml:"height"`
	Description string `xml:"description"`
}

type TextInput struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Name        string `xml:"name"`
	Link        string `xml:"link"`
}

type Item struct {
	Title       string     `xml:"rss title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Author      string     `xml:"author"`
	Categories  []Category `xml:"category"`
	Comments    string     `xml:"comments"`
	Enclosure   Enclosure  `xml:"enclosure"`
	Guid        Guid       `xml:"guid"`
	PubDate     Time       `xml:"pubDate"`
	Source      Source     `xml:"source"`
	media.Media
	dublincore.DublinCore
	content.Content
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	Text   string `xml:",chardata"`
}

type Guid struct {
	PermaLink string `xml:"permaLink,attr"`
	Text      string `xml:",chardata"`
}

type Source struct {
	Url  string `xml:"url,attr"`
	Text string `xml:",chardata"`
}

type Time struct {
	Text string `xml:",chardata"`
}

func (z Time) GetTime() time.Time {
	result, _ := parseTime(z.Text)
	return result
}

// taken from https://github.com/jteeuwen/go-pkg-rss/ timeparser.go
func parseTime(formatted string) (time.Time, error) {
	var layouts = [...]string{
		"Mon, _2 Jan 2006 15:04:05 MST",
		"Mon, _2 Jan 2006 15:04:05 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		"Mon, 2, Jan 2006 15:4",
		"02 Jan 2006 15:04:05 MST",
		"_2 Jan 2006 15:04:05 +0000", // found in the wild, differs slightly from RFC822Z
		"2006-01-02 15:04:05",        // found in the wild, apparently RFC3339 was too difficult
		"_2 Jan 2006",
		"2006-01-02",
	}
	var t time.Time
	var err error
	formatted = strings.TrimSpace(formatted)
	for _, layout := range layouts {
		t, err = time.Parse(layout, formatted)
		if !t.IsZero() {
			break
		}
	}
	return t, err
}
