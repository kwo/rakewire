package fetch

import (
	"bufio"
	"io"
	"strings"
)

// URLListToRequestArray parse url list to feeds
func URLListToRequestArray(r io.Reader) []*Request {

	var result []*Request
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			req := &Request{
				URL: url,
			}
			result = append(result, req)
		}
	}

	return result

}
