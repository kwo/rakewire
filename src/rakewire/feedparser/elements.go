package feedparser

import (
	"encoding/xml"
	"errors"
	"strings"
)

// Element encapsulates xml.Name plus attributes
type Element struct {
	name xml.Name
	attr []xml.Attr
}

// Match returns case-insensitive match of given xml.Name
func (z *Element) Match(space string, local string) bool {
	return strings.ToLower(z.name.Space) == strings.ToLower(space) && strings.ToLower(z.name.Local) == strings.ToLower(local)
}

// Attr returns the value for the given attribute
func (z *Element) Attr(space string, local string) string {
	for _, a := range z.attr {
		if strings.ToLower(a.Name.Space) == strings.ToLower(space) && strings.ToLower(a.Name.Local) == strings.ToLower(local) {
			return a.Value
		}
	}
	return ""
}

// Elements maintains a stack of Element objects
type Elements struct {
	elements []*Element
}

// Push xml StartElement on to Elements stack
func (z *Elements) Push(t xml.StartElement) *Element {
	e := &Element{name: t.Name, attr: t.Attr}
	z.elements = append(z.elements, e)
	//fmt.Printf("push: %s\n", e.Local)
	return e
}

// Pop xml EndElement off of Elements stack
func (z *Elements) Pop(t xml.EndElement) (e *Element, err error) {
	//fmt.Printf("pop:  %s\n", t.Name.Local)
	lastIndex := len(z.elements) - 1
	//fmt.Printf("lastIndex: %d\n", lastIndex)
	e = z.elements[lastIndex]
	z.elements = z.elements[:lastIndex]
	if !e.Match(t.Name.Space, t.Name.Local) {
		err = errors.New("EndElement does not match popped element")
	}
	return
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

// On returns true if the element on top matches
func (z *Elements) On(space string, local string) bool {
	e := z.peek(0)
	if e == nil {
		return false
	}
	return e.Match(space, local)
}

// In returns true if the element 1 from top matches
func (z *Elements) In(space string, local string) bool {
	e := z.peek(1)
	if e == nil {
		return false
	}
	return e.Match(space, local)
}

// Peek at the element on top of the stack
func (z *Elements) Peek() *Element {
	return z.peek(0)
}

func (z *Elements) peek(x int) *Element {
	index := len(z.elements) - x - 1
	if index < 0 || index >= len(z.elements) {
		return nil
	}
	return z.elements[index]
}
