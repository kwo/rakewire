package feedparser

type divelement struct {
	Div xmltext `xml:"div"`
}
type xmltext struct {
	Text string `xml:",innerxml"`
}

// Person atom construct
type Person struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri"`
	Email string `xml:"email"`
}

func (p *Person) String() string {
	var result string
	if p != nil {
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
