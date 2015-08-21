package feedparser

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Element encapsulates xml.Name plus attributes
type Element struct {
	name xml.Name
	attr []xml.Attr
}

// Match returns case-insensitive match of given xml.Name
func (z *Element) Match(space string, local string) bool {
	return strings.ToLower(z.name.Local) == strings.ToLower(local) && strings.ToLower(z.name.Space) == strings.ToLower(space)
}

// Attr returns the value for the given attribute
func (z *Element) Attr(space string, local string) string {
	for _, a := range z.attr {
		if strings.ToLower(a.Name.Local) == strings.ToLower(local) && strings.ToLower(a.Name.Space) == strings.ToLower(space) {
			return a.Value
		}
	}
	return ""
}

// Elements maintains a stack of Element objects
type Elements struct {
	elements []*Element
}

// IsStackFeed if you are at the feed level
func (z *Elements) IsStackFeed(args ...int) bool {

	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}

	length := z.Level() - offset

	switch length {

	case 1:
		for i := 0; i < length; i++ {
			e := z.elements[i]
			switch i {
			case 0:
				if e.name.Space != nsAtom || e.name.Local != "feed" {
					return false
				}
			default:
				return false
			} // switch
		} // loop
		return true

	case 2:
		for i := 0; i < length; i++ {
			e := z.elements[i]
			switch i {
			case 0:
				if e.name.Space != nsRSS || e.name.Local != "rss" {
					return false
				}
			case 1:
				if e.name.Space != nsRSS || e.name.Local != "channel" {
					return false
				}
			default:
				return false
			} // switch
		} // loop
		return true
	} // level switch

	return false

}

// IsStackEntry if you are at the entry level
func (z *Elements) IsStackEntry(args ...int) bool {

	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}

	length := z.Level() - offset

	switch length {

	case 2:
		for i := 0; i < length; i++ {
			e := z.elements[i]
			switch i {
			case 0:
				if e.name.Space != nsAtom || e.name.Local != "feed" {
					return false
				}
			case 1:
				if e.name.Space != nsAtom || e.name.Local != "entry" {
					return false
				}
			default:
				return false
			} // switch
		} // loop
		return true

	case 3:
		for i := 0; i < length; i++ {
			e := z.elements[i]
			switch i {
			case 0:
				if e.name.Space != nsRSS || e.name.Local != "rss" {
					return false
				}
			case 1:
				if e.name.Space != nsRSS || e.name.Local != "channel" {
					return false
				}
			case 2:
				if e.name.Space != nsRSS || e.name.Local != "item" {
					return false
				}
			default:
				return false
			} // switch
		} // loop
		return true
	} // level switch

	return false

}

// Attr walks down the stack delivering the first matching attribute value
func (z *Elements) Attr(space string, local string) string {
	for i := len(z.elements) - 1; i >= 0; i-- {
		if value := z.elements[i].Attr(space, local); value != "" {
			return value
		}
	}
	return ""
}

// Level returns the depth of the stack
func (z *Elements) Level() int {
	return len(z.elements)
}

// Peek at the element on top of the stack
func (z *Elements) Peek() *Element {
	return z.peek(0)
}

// PeekIf at the element on top of the stack if match
func (z *Elements) PeekIf(t xml.EndElement) (*Element, error) {
	e := z.Peek()
	if e.Match(t.Name.Space, t.Name.Local) {
		return e, nil
	}
	return nil, fmt.Errorf("%s:%s does not match %s:%s", e.name.Space, e.name.Local, t.Name.Space, t.Name.Local)
}

// Pop xml EndElement off of Elements stack
func (z *Elements) Pop() *Element {
	lastIndex := len(z.elements) - 1
	e := z.elements[lastIndex]
	z.elements = z.elements[:lastIndex]
	return e
}

// PopIf pop xml EndElement off of Elements stack if match
func (z *Elements) PopIf(t xml.EndElement) (*Element, error) {
	e := z.Peek()
	if e.Match(t.Name.Space, t.Name.Local) {
		z.Pop()
		return e, nil
	}
	return nil, fmt.Errorf("%s:%s does not match %s:%s", e.name.Space, e.name.Local, t.Name.Space, t.Name.Local)
}

// Push xml StartElement on to Elements stack
func (z *Elements) Push(t xml.StartElement) *Element {
	e := &Element{name: t.Name, attr: t.Attr}
	z.elements = append(z.elements, e)
	return e
}

// String prints the stack
func (z *Elements) String() string {
	var result []string
	for _, e := range z.elements {
		result = append(result, e.name.Local)
	}
	return strings.Join(result, ">")
}

func (z *Elements) peek(x int) *Element {
	index := len(z.elements) - x - 1
	if index < 0 || index >= len(z.elements) {
		return nil
	}
	return z.elements[index]
}
