package opml

import (
	"encoding/xml"
	"io"
	"strings"
	"time"
)

// OPML represents the top-level opml document
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Head    *Head
	Body    *Body
}

// Head holds some meta information about the document.
type Head struct {
	Title        string    `xml:"title"`
	DateCreated  time.Time `xml:"dateCreated,omitempty"`
	DateModified string    `xml:"dateModified,omitempty"`
	OwnerName    string    `xml:"ownerName,omitempty"`
	OwnerEmail   string    `xml:"ownerEmail,omitempty"`
	OwnerID      string    `xml:"ownerId,omitempty"`
}

// Body represents the main opml document body
type Body struct {
	XMLName  xml.Name   `xml:"body"`
	Outlines []*Outline `xml:"outline"`
}

// Outline holds all information about an outline.
type Outline struct {
	Language    string     `xml:"language,attr,omitempty"`
	Version     string     `xml:"version,attr,omitempty"`
	Created     *time.Time `xml:"created,attr,omitempty"`
	Type        string     `xml:"type,attr,omitempty"`
	Category    string     `xml:"category,attr,omitempty"`
	Text        string     `xml:"text,attr"`
	Title       string     `xml:"title,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`
	XMLURL      string     `xml:"xmlUrl,attr,omitempty"`
	HTMLURL     string     `xml:"htmlUrl,attr,omitempty"`
	Outlines    []*Outline `xml:"outline"`
}

// Branch is used to group outlines when flattening
type Branch struct {
	Name     string
	AutoRead bool
	AutoStar bool
}

type container interface {
	getOutlines() []*Outline
}

func (z *Body) getOutlines() []*Outline {
	return z.Outlines
}

func (z *Outline) getOutlines() []*Outline {
	return z.Outlines
}

// Parse parses input into a OPML structure
func Parse(r io.Reader) (*OPML, error) {
	o := &OPML{}
	err := xml.NewDecoder(r).Decode(o)
	return o, err
}

// Format serializes output to a writer
func Format(o *OPML, w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	return encoder.Encode(o)
}

// Flatten pulls nested groups up to the top level, modifing the group name.
func Flatten(container container) map[*Branch][]*Outline {
	result := make(map[*Branch][]*Outline)
	flattenContainer(container, &Branch{Name: ""}, result)
	return result
}

func flattenContainer(container container, branch *Branch, result map[*Branch][]*Outline) {
	for _, outline := range container.getOutlines() {
		if outline.Type == "rss" {
			result[branch] = append(result[branch], outline)
		} else {
			// inherited from parent branch
			b := &Branch{
				AutoRead: branch.AutoRead || strings.Contains(outline.Category, "autoread"),
				AutoStar: branch.AutoStar || strings.Contains(outline.Category, "autostar"),
			}
			if branch.Name == "" {
				b.Name = outline.Text
				flattenContainer(outline, b, result) // 2nd level
			} else {
				b.Name = branch.Name + "/" + outline.Text
				flattenContainer(outline, b, result) // 3rd + levels only
			}
		}
	}
}
