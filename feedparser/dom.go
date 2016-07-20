package feedparser

import (
	"strings"
)

type content struct {
	Type  string  `xml:"type,attr"`
	Text  string  `xml:",chardata"`
	XHtml xmlData `xml:"div"`
}

type xmlData struct {
	Text string `xml:",innerxml"`
}

func (z *content) ToString() string {
	result := strings.TrimSpace(z.XHtml.Text)
	if isEmpty(result) {
		result = strings.TrimSpace(z.Text)
	}
	return result
}

type generator struct {
	Text    string `xml:",chardata"`
	URI     string `xml:"uri,attr"`
	Version string `xml:"version,attr"`
}

func (z *generator) ToString() string {
	z.Text = strings.TrimSpace(z.Text)
	z.URI = strings.TrimSpace(z.URI)
	z.Version = strings.TrimSpace(z.Version)
	result := z.Text
	if !isEmpty(result) {
		if !isEmpty(z.Version) {
			result += " " + z.Version
		}
		if !isEmpty(z.URI) {
			result += " (" + z.URI + ")"
		}
	}
	return result
}

type rssImage struct {
	URL string `xml:"url"`
}

type person struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

func (p *person) ToString() string {
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.TrimSpace(p.Email)
	p.URI = strings.TrimSpace(p.URI)
	result := p.Name
	if !isEmpty(p.Email) {
		result += " <" + p.Email + ">"
	}
	if !isEmpty(p.URI) {
		result += " (" + p.URI + ")"
	}
	return result
}

type text struct {
	Text string `xml:",chardata"`
}

func (z *text) ToString() string {
	return strings.TrimSpace(z.Text)
}
