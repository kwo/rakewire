package feedparser

import (
	"strings"
	"testing"
)

func TestRewriteSimpleHref(t *testing.T) {

	content := `
	<p>abc <a href="/123" title="hello">hello</a> def <img src="/image"/> ghi <img src="/pict.jpg"></img></p>
	`
	content2 := `
	<p>abc <a href="http://localhost/123" title="hello">hello</a> def <img src="http://localhost/image"/> ghi <img src="http://localhost/pict.jpg"></img></p>
	`
	result := RewriteContentWithAbsoluteURLs("http://localhost", content)

	t.Log(result)

	if strings.Compare(strings.TrimSpace(result), strings.TrimSpace(content2)) != 0 {
		t.Errorf("Failed to make URLs absolute: %s", result)
	}

}

func TestRewriteOngoing(t *testing.T) {

	f := testFile(t, "testdata/ongoing.atom")
	if f == nil {
		t.Fatal("Cannot parse feed.")
	}
	RewriteFeedWithAbsoluteURLs(f)

	for _, entry := range f.Entries {
		t.Logf("entry: %s", entry.Title)
		if entry.Title == "New British Isles" {
			links := FindURLs(entry.Content)
			if len(links) != 5 {
				t.Fatalf("bad link count %d, expecting %d", len(links), 5)
			}
			expectedLink := "https://www.tbray.org/ongoing/When/200x/2004/11/04/NewNorthAmerica"
			if links[4].URL.String() != expectedLink {
				t.Errorf("bad link: %s, expecting %s", links[4].URL, expectedLink)
			}

		}
	}

}
