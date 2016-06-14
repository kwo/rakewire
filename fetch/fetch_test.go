package fetch

import (
	"github.com/kwo/rakewire/model"
	"net/url"
	"testing"
)

func TestInterfaceService(t *testing.T) {

	var s model.Service = &Service{}
	if s == nil {
		t.Fatal("Does not implement m.Service interface.")
	}

}

func TestURL1(t *testing.T) {

	u1 := "http://www.recode.net/blog.xml"
	u2 := "/rss/index.xml"
	u3 := "http://www.recode.net/rss/index.xml"

	url1, errParse1 := url.Parse(u1)
	if errParse1 != nil {
		t.Fatalf("cannot parse url: %s. %s", u1, errParse1.Error())
	}

	url2, errParse2 := url.Parse(u2)
	if errParse2 != nil {
		t.Fatalf("cannot parse url: %s. %s", u2, errParse2.Error())
	}

	url3 := url1.ResolveReference(url2)
	if url3.String() != u3 {
		t.Errorf("Bad URL: %s, expected %s", url3.String(), u3)
	}

}

func TestURL2(t *testing.T) {

	u1 := "http://www.recode.net/blog.xml"
	u2 := "rss/index.xml"
	u3 := "http://www.recode.net/rss/index.xml"

	url1, errParse1 := url.Parse(u1)
	if errParse1 != nil {
		t.Fatalf("cannot parse url: %s. %s", u1, errParse1.Error())
	}

	url2, errParse2 := url.Parse(u2)
	if errParse2 != nil {
		t.Fatalf("cannot parse url: %s. %s", u2, errParse2.Error())
	}

	url3 := url1.ResolveReference(url2)
	if url3.String() != u3 {
		t.Errorf("Bad URL: %s, expected %s", url3.String(), u3)
	}

}
