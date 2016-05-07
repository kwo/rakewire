package opml

import (
	"bufio"
	"bytes"
	"io"
	"github.com/kwo/rakewire/model"
	"strings"
	"testing"
)

func TestExport(t *testing.T) {

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)
	user := addUser(t, db)

	importData(t, db, user, getImport1())

	exportBuffer := &bytes.Buffer{}
	exportData(t, db, user, exportBuffer)

	actualLines := readLines(t, exportBuffer)
	expectedLines := readLines(t, getExport1())

	actualLineCount := len(actualLines)
	expectedLineCount := len(expectedLines)

	if actualLineCount != expectedLineCount {
		t.Log("Actual")
		for n, line := range actualLines {
			t.Logf("%02d %s", n, line)
		}
		t.Log("Expected")
		for n, line := range expectedLines {
			t.Logf("%02d %s", n, line)
		}
		t.Fatalf("Bad line count: %d, expected: %d", actualLineCount, expectedLineCount)
	}

	for n := 0; n < actualLineCount; n++ {
		aLine := strings.SplitAfter(actualLines[n], " created=")[0]
		eLine := strings.SplitAfter(expectedLines[n], " created=")[0]
		if !strings.Contains(aLine, "<dateCreated>") {
			if aLine != eLine {
				t.Errorf("line mismatch: %s, expected %s", aLine, eLine)
			}
		}
	}

}

func readLines(t *testing.T, reader io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	err := scanner.Err()
	if err != nil {
		t.Fatalf("Error reading lines: %s", err.Error())
	}
	return lines
}

func exportData(t *testing.T, db model.Database, user *model.User, writer io.Writer) {

	err := db.Select(func(tx model.Transaction) error {

		// import
		opml, err := Export(tx, user)
		if err != nil {
			return err
		}

		return Format(opml, writer)

	})
	if err != nil {
		t.Errorf("Error exporting OPML: %s", err.Error())
	}

}

func getExport1() io.Reader {
	document := `
	<?xml version="1.0" encoding="UTF-8"?>
	<opml>
	  <head>
	    <title>Rakewire Subscriptions</title>
			<dateCreated>2016-03-16T21:04:10Z</dateCreated>
	    <ownerName>opml.ninja</ownerName>
	  </head>
	  <body>
	    <outline title="Group1">
	      <outline type="rss" title="g1title1" xmlUrl="g1xmlurl1" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g1title2" xmlUrl="g1xmlurl2" created="2016-03-16T20:57:33Z" category="+autostar"></outline>
	      <outline type="rss" title="g1title3" xmlUrl="g1xmlurl3" created="2016-03-16T20:57:33Z"></outline>
	    </outline>
	    <outline title="Group2">
	      <outline type="rss" title="g2title1" xmlUrl="g2xmlurl1" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g2title2" xmlUrl="g2xmlurl2" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g2title3" xmlUrl="g2xmlurl3" created="2016-03-16T20:57:33Z"></outline>
	    </outline>
	    <outline title="GroupX/Group3">
	      <outline type="rss" title="g3title1" xmlUrl="g3xmlurl1" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g3title2" xmlUrl="g3xmlurl2" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g3title3" xmlUrl="g3xmlurl3" created="2016-03-16T20:57:33Z"></outline>
	    </outline>
	    <outline title="GroupX/Group4">
	      <outline type="rss" title="g4title1" xmlUrl="g4xmlurl1" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g4title2" xmlUrl="g4xmlurl2" created="2016-03-16T20:57:33Z"></outline>
	      <outline type="rss" title="g4title3" xmlUrl="g4xmlurl3" created="2016-03-16T20:57:33Z"></outline>
	    </outline>
	  </body>
	</opml>`
	return strings.NewReader(document)
}
