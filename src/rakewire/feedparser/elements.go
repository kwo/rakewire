package feedparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// Element encapsulates xml.Name plus attributes
type Element struct {
	xml.Name
	Attr map[xml.Name]string
}

// Match returns case-insensitive match of given xml.Name
func (z *Element) Match(space string, local string) bool {
	return strings.ToLower(z.Space) == strings.ToLower(space) && strings.ToLower(z.Local) == strings.ToLower(local)
}

func elementAttr(attrs map[xml.Name]string, space string, local string) string {
	return attrs[xml.Name{Space: space, Local: local}]
}

func elementLowerName(n xml.Name) xml.Name {
	return xml.Name{
		Local: strings.ToLower(n.Local),
		Space: strings.ToLower(n.Space),
	}
}

func elementMapAttr(attrs []xml.Attr) map[xml.Name]string {
	result := make(map[xml.Name]string)
	for _, attr := range attrs {
		result[elementLowerName(attr.Name)] = attr.Value
	}
	return result
}

// Elements maintains a stack of Element objects
type Elements struct {
	elements []*Element
}

// Push xml StartElement on to Elements stack
func (z *Elements) Push(t xml.StartElement) *Element {
	e := &Element{Name: t.Name, Attr: elementMapAttr(t.Attr)}
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

// Within returns true if the specified name is contained within the stack
func (z *Elements) Within(space string, local string) bool {
	if len(z.elements) == 0 {
		return false
	}
	for i := 0; i < len(z.elements); i++ {
		e := z.peek(i)
		fmt.Printf("i: %s:%s %d %d %t \n", space, local, i, len(z.elements), e == nil)
		if e.Match(space, local) {
			return true
		}
	}
	return false
}

func (z *Elements) peek(x int) *Element {
	index := len(z.elements) - x - 1
	if index < 0 || index >= len(z.elements) {
		return nil
	}
	return z.elements[index]
}
