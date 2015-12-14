package fever

import (
	"encoding/xml"
	"fmt"
	m "rakewire/model"
	"strconv"
	"time"
)

// API top level struct
type API struct {
	prefix string
	db     Database
}

// Database defines the interface to the database
type Database interface {
	UserGetByFeverHash(feverhash string) (*m.User, error)
}

// Response defines the json/xml response return by requests.
type Response struct {
	XMLName       xml.Name  `json:"-" xml:"response"`
	Version       int       `json:"api_version" xml:"api_version"`
	Authorized    int       `json:"auth" xml:"auth"`
	LastRefreshed feverTime `json:"last_refreshed_on_time,omitempty" xml:"last_refreshed_on_time,omitempty"`
}

type feverTime struct {
	time.Time
}

func (z feverTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", strconv.FormatInt(z.Unix(), 10))), nil
}

func (z feverTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !z.IsZero() {
		e.EncodeToken(start)
		e.EncodeToken(xml.CharData([]byte(strconv.FormatInt(z.Unix(), 10))))
		e.EncodeToken(xml.EndElement{Name: start.Name})
	}
	return nil
}
