package feedparser

import (
	"strings"
)

type person struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

func (p *person) String() string {
	var result string
	if p != nil {
		p.Name = strings.TrimSpace(p.Name)
		p.Email = strings.TrimSpace(p.Email)
		p.URI = strings.TrimSpace(p.URI)
		result = p.Name
		if p.Email != "" {
			result += " <" + p.Email + ">"
		}
		if p.URI != "" {
			result += " (" + p.URI + ")"
		}
	}
	return result
}

type xmlText struct {
	Type  string  `xml:"type,attr"`
	Text  string  `xml:",chardata"`
	XHtml xmlData `xml:"div"`
}

type xmlData struct {
	Text string `xml:",innerxml"`
}

func (z *xmlText) GetText() *Text {
	var result Text
	// #TODO:0 use base to fix relative HREFs in XML
	result.Type = z.Type
	if result.Type == "" {
		result.Type = "text"
	}
	result.Text = strings.TrimSpace(z.XHtml.Text)
	if result.Text == "" {
		result.Text = strings.TrimSpace(z.Text)
	}
	return &result
}

type xmlString struct {
	Text string `xml:",chardata"`
}
