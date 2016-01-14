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
	XMLName      xml.Name  `xml:"head"`
	Title        string    `xml:"title"`
	DateCreated  time.Time `xml:"dateCreated,omitempty"`
	DateModified string    `xml:"dateModified,omitempty"`
	OwnerName    string    `xml:"ownerName,omitempty"`
	OwnerEmail   string    `xml:"ownerEmail,omitempty"`
	OwnerID      string    `xml:"ownerId,omitempty"`
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
	Text        string     `xml:"text,attr"`
	Title       string     `xml:"title,attr,omitempty"`
	XMLURL      string     `xml:"xmlUrl,attr,omitempty"`
	HTMLURL     string     `xml:"htmlUrl,attr,omitempty"`
	Category    string     `xml:"category,attr,omitempty"`
	Created     *time.Time `xml:"created,attr,omitempty"`
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

// Sort sort by title
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

// Flatten pulls nested groups up to the top level, modifing the group name.
func Flatten(outlines Outlines) map[*Outline]Outlines {
	result := make(map[*Outline]Outlines)
	flatten(outlines, &Outline{}, result)
	return result
}

func flatten(outlines Outlines, branch *Outline, result map[*Outline]Outlines) {
	for _, outline := range outlines {
		if outline.Type == "rss" {
			result[branch] = append(result[branch], outline)
		} else {
			// inherited from parent branch
			b := &Outline{}
			b.SetAutoRead(branch.IsAutoRead() || outline.IsAutoRead())
			b.SetAutoStar(branch.IsAutoStar() || outline.IsAutoStar())
			if branch.Text == "" {
				b.Text = outline.Text
				b.Title = outline.Title
				flatten(outline.Outlines, b, result) // 2nd level
			} else {
				b.Text = branch.Text + "/" + outline.Text
				b.Title = branch.Title + "/" + outline.Title
				flatten(outline.Outlines, b, result) // 3rd + levels only
			}
		}
	}
}

func groupOutlinesByURL(flatOPML map[*Outline]Outlines) map[string]*Outline {
	result := make(map[string]*Outline)
	for _, outlines := range flatOPML {
		for _, outline := range outlines {
			if _, ok := result[outline.XMLURL]; !ok {
				result[outline.XMLURL] = outline
			}
		}
	}
	return result
}
