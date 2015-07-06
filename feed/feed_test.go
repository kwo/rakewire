package feed

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestFeed(t *testing.T) {
	Parse("https://ostendorf.com/feed.xml")
}
