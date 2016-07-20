package feedparser

import (
	"bytes"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

// Link represents an HTML link
type Link struct {
	Text string
	URL  *url.URL
}

func findURLs(content string) []*Link {

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

// base url should be the alternate link to the feed entry
func makeAbsoluteURLs(base, content string) string {

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
