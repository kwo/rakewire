package feedparser

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an HTML link
type Link struct {
	Text string
	URL  *url.URL
}

// FindURLs will find all HTML links in the given content.
func FindURLs(content string) []*Link {

	if isEmpty(content) {
		return nil
	}

	result := []*Link{}
	tokenizer := html.NewTokenizer(strings.NewReader(content))
	var link *Link

Loop:
	for {

		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			break Loop
		}

		token := tokenizer.Token()

		if tt == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					linkURL, err := url.Parse(attr.Val)
					if err == nil {
						link = &Link{URL: linkURL}
					}
					break // exit attr loop
				}
			}
		} else if tt == html.EndTagToken && token.Data == "a" && link != nil {
			link.Text = strings.TrimSpace(link.Text)
			result = append(result, link)
			link = nil
		} else if link != nil {
			text := strings.Join(strings.Fields(token.String()), " ")
			link.Text = link.Text + text + " "
		}

	} // loop

	return result

}

// RewriteFeedWithAbsoluteURLs will rewrite all feed entries with absolute URLs if not already.
func RewriteFeedWithAbsoluteURLs(f *Feed) {
	for _, entry := range f.Entries {
		RewriteEntryWithAbsoluteURLs(entry)
	}
}

// RewriteEntryWithAbsoluteURLs will rewrite all links in Entry.Content and Entry.Summary as absolute URLs if not already.
func RewriteEntryWithAbsoluteURLs(entry *Entry) {
	entry.Content = RewriteContentWithAbsoluteURLs(entry.LinkAlternate, entry.Content)
	entry.Summary = RewriteContentWithAbsoluteURLs(entry.LinkAlternate, entry.Summary)
}

// RewriteContentWithAbsoluteURLs will rewrite all HTML anchor and image tags in content as absolute URLs,
// if not already, using the given base url as a reference.
func RewriteContentWithAbsoluteURLs(base, content string) string {

	if isEmpty(base) || isEmpty(content) {
		return content
	}

	baseURL, errBase := url.Parse(base)
	if errBase != nil {
		return content
	}

	result := &bytes.Buffer{}
	tokenizer := html.NewTokenizer(strings.NewReader(content))

Loop:
	for {

		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			break Loop
		}

		token := tokenizer.Token()
		if (tt == html.StartTagToken || tt == html.SelfClosingTagToken) && (token.Data == "a" || token.Data == "img") {
			for index, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					token.Attr[index] = html.Attribute{Namespace: attr.Namespace, Key: attr.Key, Val: rewriteURL(baseURL, attr.Val)}
				}
			}
		}

		result.WriteString(token.String())

	} // loop

	return string(result.Bytes())

}

func rewriteURL(baseURL *url.URL, ref string) string {

	refURL, errRef := url.Parse(ref)
	if errRef != nil {
		return ref
	}

	if refURL.IsAbs() {
		return ref
	}

	return baseURL.ResolveReference(refURL).String()

}
