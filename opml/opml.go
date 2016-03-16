package opml

import (
	"encoding/xml"
	"io"
	"sort"
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
	XMLName      xml.Name   `xml:"head"`
	Title        string     `xml:"title"`
	DateCreated  *time.Time `xml:"dateCreated,omitempty"`
	DateModified string     `xml:"dateModified,omitempty"`
	OwnerName    string     `xml:"ownerName,omitempty"`
	OwnerEmail   string     `xml:"ownerEmail,omitempty"`
	OwnerID      string     `xml:"ownerId,omitempty"`
}

// Body represents the main opml document body
type Body struct {
	XMLName  xml.Name `xml:"body"`
	Outlines Outlines `xml:"outline"`
}

// Outline holds all information about an outline.
type Outline struct {
	Language    string     `xml:"language,attr,omitempty"`
	Version     string     `xml:"version,attr,omitempty"`
	Type        string     `xml:"type,attr,omitempty"`
	Title       string     `xml:"title,attr"`
	Text        string     `xml:"text,attr,omitempty"` // always mirrors title
	XMLURL      string     `xml:"xmlUrl,attr,omitempty"`
	HTMLURL     string     `xml:"htmlUrl,attr,omitempty"`
	Created     *time.Time `xml:"created,attr,omitempty"`
	Category    string     `xml:"category,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`
	Outlines    Outlines   `xml:"outline"`
}

// IsAutoRead returns if the outline is marked for autoread
func (z *Outline) IsAutoRead() bool {
	return strings.Contains(z.Category, "+autoread")
}

// IsAutoStar returns if the outline is marked for autostar
func (z *Outline) IsAutoStar() bool {
	return strings.Contains(z.Category, "+autostar")
}

// SetAutoRead sets the autoread flag
func (z *Outline) SetAutoRead(value bool) {
	if value && !z.IsAutoRead() {
		z.Category += " +autoread"
	} else if !value && z.IsAutoRead() {
		z.Category = strings.Replace(z.Category, "+autoread", "", -1)
	}
	z.Category = strings.TrimSpace(z.Category)
}

// SetAutoStar sets the autostar flag
func (z *Outline) SetAutoStar(value bool) {
	if value && !z.IsAutoStar() {
		z.Category += " +autostar"
	} else if !value && z.IsAutoStar() {
		z.Category = strings.Replace(z.Category, "+autostar", "", -1)
	}
	z.Category = strings.TrimSpace(z.Category)
}

// Outlines array type
type Outlines []*Outline

func (z Outlines) Len() int      { return len(z) }
func (z Outlines) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Outlines) Less(i, j int) bool {
	return strings.ToLower(z[i].Title) < strings.ToLower(z[j].Title)
}

// Sort recursively sort by title
func (z Outlines) Sort() {
	sort.Sort(z)
	for _, outline := range z {
		outline.Outlines.Sort()
	}
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
