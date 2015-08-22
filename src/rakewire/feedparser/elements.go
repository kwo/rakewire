package feedparser

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Element encapsulates xml.Name plus attributes
type element struct {
	name xml.Name
	attr []xml.Attr
}

// Match returns case-insensitive match of given xml.Name
func (z *element) Match(space string, local string) bool {
	return strings.ToLower(z.name.Local) == strings.ToLower(local) && strings.ToLower(z.name.Space) == strings.ToLower(space)
}

// Attr returns the value for the given attribute
func (z *element) Attr(space string, local string) string {
	for _, a := range z.attr {
		if strings.ToLower(a.Name.Local) == strings.ToLower(local) && strings.ToLower(a.Name.Space) == strings.ToLower(space) {
			return a.Value
		}
	}
	return ""
}

// elements maintains a stack of Element objects
type elements struct {
	stack []*element
}

// IsStackFeed if you are at the feed level
func (z *elements) IsStackFeed(args ...int) bool {

	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}

	length := z.Level() - offset

	switch length {

	case 1:
		for i := 0; i < length; i++ {
			e := z.stack[i]
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
			e := z.stack[i]
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
func (z *elements) IsStackEntry(args ...int) bool {

	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}

	length := z.Level() - offset

	switch length {

	case 2:
		for i := 0; i < length; i++ {
			e := z.stack[i]
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
			e := z.stack[i]
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
func (z *elements) Attr(space string, local string) string {
	for i := len(z.stack) - 1; i >= 0; i-- {
		if value := z.stack[i].Attr(space, local); value != "" {
			return value
		}
	}
	return ""
}

// Level returns the depth of the stack
func (z *elements) Level() int {
	return len(z.stack)
}

// Peek at the element on top of the stack
func (z *elements) Peek() *element {
	return z.peek(0)
}

// PeekIf at the element on top of the stack if match
func (z *elements) PeekIf(t xml.EndElement) (*element, error) {
	e := z.Peek()
	if e.Match(t.Name.Space, t.Name.Local) {
		return e, nil
	}
	return nil, fmt.Errorf("%s:%s does not match %s:%s", e.name.Space, e.name.Local, t.Name.Space, t.Name.Local)
}

// Pop xml EndElement off of elements stack
func (z *elements) Pop() *element {
	lastIndex := len(z.stack) - 1
	e := z.stack[lastIndex]
	z.stack = z.stack[:lastIndex]
	return e
}

// PopIf pop xml EndElement off of elements stack if match
func (z *elements) PopIf(t xml.EndElement) (*element, error) {
	e := z.Peek()
	if e.Match(t.Name.Space, t.Name.Local) {
		z.Pop()
		return e, nil
	}
	return nil, fmt.Errorf("%s:%s does not match %s:%s", e.name.Space, e.name.Local, t.Name.Space, t.Name.Local)
}

// Push xml StartElement on to elements stack
func (z *elements) Push(t xml.StartElement) *element {
	e := &element{name: t.Name, attr: t.Attr}
	z.stack = append(z.stack, e)
	return e
}

// String prints the stack
func (z *elements) String() string {
	var result []string
	for _, e := range z.stack {
		result = append(result, e.name.Local)
	}
	return strings.Join(result, ">")
}

func (z *elements) peek(x int) *element {
	index := len(z.stack) - x - 1
	if index < 0 || index >= len(z.stack) {
		return nil
	}
	return z.stack[index]
}
