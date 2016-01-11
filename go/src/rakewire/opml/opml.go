package opml

import (
	"encoding/xml"
	"io"
)

// OPML represents the top-level opml document
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Body    *Body
}

// Body represents the main opml document body
type Body struct {
	XMLName  xml.Name   `xml:"body"`
	Outlines []*Outline `xml:"outline"`
}

// GetOutlines returns outlines in the body
func (z *Body) GetOutlines() []*Outline {
	return z.Outlines
}

// Outline holds all information about an outline.
type Outline struct {
	Language    string     `xml:"language,attr,omitempty"`
	Version     string     `xml:"version,attr,omitempty"`
	Created     string     `xml:"created,attr,omitempty"`
	Type        string     `xml:"type,attr,omitempty"`
	Category    string     `xml:"category,attr,omitempty"`
	Text        string     `xml:"text,attr"`
	Title       string     `xml:"title,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`
	XMLURL      string     `xml:"xmlUrl,attr,omitempty"`
	HTMLURL     string     `xml:"htmlUrl,attr,omitempty"`
	URL         string     `xml:"url,attr,omitempty"`
	Outlines    []*Outline `xml:"outline"`
}

// GetOutlines returns nested outlines in the outline
func (z *Outline) GetOutlines() []*Outline {
	return z.Outlines
}

// Container contains outlines
type Container interface {
	GetOutlines() []*Outline
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

// Flatten groups outlines by parent path
func Flatten(container Container) map[string][]*Outline {
	result := make(map[string][]*Outline)
	flatten(container, "", result)
	return result
}

func flatten(container Container, path string, result map[string][]*Outline) {
	for _, outline := range container.GetOutlines() {
		if outline.XMLURL != "" {
			result[path] = append(result[path], outline)
		} else if outline.Text != "" {
			if path != "" {
				flatten(outline, path+":"+outline.Text, result)
			} else {
				flatten(outline, outline.Text, result)
			}
		} else {
			flatten(outline, path, result)
		}
	}
}
