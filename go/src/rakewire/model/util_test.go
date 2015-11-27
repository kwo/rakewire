package model

import (
	"testing"
)

func TestParseFeedsFromFile(t *testing.T) {

	t.Parallel()

	feeds, err := ParseFeedsFromFile("../../../test/feedlistmini.txt")
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 10, len(feeds))
	assertEqual(t, "http://www.addrup.de/feed.xml", feeds[0].URL)

}

func TestParseFeedsFromFileError(t *testing.T) {

	t.Parallel()

	_, err := ParseFeedsFromFile("../../../test/feedlistmini2.txt")
	assertNotNil(t, err)

}
