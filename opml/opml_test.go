package opml

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestParseFormat(t *testing.T) {

	opml1 := getTestOPML()
	document := getTestDocument()

	// OPMLFormat opml1 to buffer
	buf1 := &bytes.Buffer{}
	if err := Format(opml1, buf1); err != nil {
		t.Fatalf("Error OPMLFormatting OPML1: %s", err.Error())
	}

	// OPMLParse and reformat opml2 to buffer
	opml2, err := Parse(document)
	if err != nil {
		t.Fatalf("Cannot OPMLParse opml2: %s", err.Error())
	}
	buf2 := &bytes.Buffer{}
	if err := Format(opml2, buf2); err != nil {
		t.Fatalf("Error OPMLFormatting OPML2: %s", err.Error())
	}

	// compare buffers
	if bytes.Compare(buf1.Bytes(), buf2.Bytes()) != 0 {
		t.Errorf("Format mismatch:\nExpected\n\n %s\n\nActual\n\n%s", buf2.Bytes(), buf1.Bytes())
	}

}

func TestSort(t *testing.T) {

	opml1 := getTestOPML()
	document := getTestDocumentSorted()

	opml1.Body.Outlines.Sort()

	// OPMLFormat opml1 to buffer
	buf1 := &bytes.Buffer{}
	if err := Format(opml1, buf1); err != nil {
		t.Fatalf("Error OPMLFormatting OPML1: %s", err.Error())
	}

	// OPMLParse and reformat opml2 to buffer
	opml2, err := Parse(document)
	if err != nil {
		t.Fatalf("Cannot OPMLParse opml2: %s", err.Error())
	}
	buf2 := &bytes.Buffer{}
	if err := Format(opml2, buf2); err != nil {
		t.Fatalf("Error OPMLFormatting OPML2: %s", err.Error())
	}

	// compare buffers
	if bytes.Compare(buf1.Bytes(), buf2.Bytes()) != 0 {
		t.Errorf("Format mismatch:\nExpected\n\n %s\n\nActual\n\n%s", buf2.Bytes(), buf1.Bytes())
	}

}

func TestFlatten1(t *testing.T) {

	// given
	opml1 := getTestOPML()
	flatOPML := flatten(opml1.Body.Outlines)

	// when
	expectedKeyCount := 4
	actualKeyCount := len(flatOPML)

	// then
	if actualKeyCount != expectedKeyCount {
		t.Errorf("Expected %d, actual %d", expectedKeyCount, actualKeyCount)
	}

}

func TestFlatten2(t *testing.T) {

	// given
	opml1 := getTestOPML()

	// when
	opml1.Body.Outlines[1].SetAutoRead(true) // group1
	opml1.Body.Outlines[2].SetAutoRead(true) // group3
	flatOPML := flatten(opml1.Body.Outlines)

	for group := range flatOPML {
		t.Logf("group: %t %t %s", group.IsAutoRead(), group.IsAutoStar(), group.Title)
	}

}

func TestAutoRead1(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoRead(true)

	// then
	if !outline.IsAutoRead() {
		t.Error("Expected autoread to be set but it is false")
	}
	if !strings.Contains(outline.Category, "+autoread") {
		t.Error("Expected category to contain +autoread")
	}

}

func TestAutoRead2(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoRead(true)
	outline.SetAutoRead(false)

	// then
	if outline.IsAutoRead() {
		t.Error("Expected autoread to be off but it is true")
	}
	if outline.Category != "" {
		t.Errorf("Expected category to be blank, actual %s", outline.Category)
	}

}

func TestAutoStar1(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoStar(true)

	// then
	if !outline.IsAutoStar() {
		t.Error("Expected autostar to be set but it is false")
	}
	if !strings.Contains(outline.Category, "+autostar") {
		t.Error("Expected category to contain +autostar")
	}

}

func TestAutoStar2(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoStar(true)
	outline.SetAutoStar(false)

	// then
	if outline.IsAutoStar() {
		t.Error("Expected autostar to be off but it is true")
	}
	if outline.Category != "" {
		t.Errorf("Expected category to be blank, actual %s", outline.Category)
	}

}

func TestAutoFlags1(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoRead(true)
	outline.SetAutoStar(true)

	// then
	if !outline.IsAutoRead() {
		t.Error("Expected autoread to be set but it is false")
	}
	if !strings.Contains(outline.Category, "+autoread") {
		t.Error("Expected category to contain +autoread")
	}
	if !outline.IsAutoStar() {
		t.Error("Expected autostar to be set but it is false")
	}
	if !strings.Contains(outline.Category, "+autostar") {
		t.Error("Expected category to contain +autostar")
	}

}

func TestAutoFlags2(t *testing.T) {

	// given
	outline := &Outline{}

	// when
	outline.SetAutoStar(true)
	outline.SetAutoRead(true)
	outline.SetAutoRead(false)

	// then
	if outline.IsAutoRead() {
		t.Error("Expected autoread to be off but it is true")
	}
	if !outline.IsAutoStar() {
		t.Error("Expected autostar to be set but it is false")
	}
	if outline.Category != "+autostar" {
		t.Errorf("Expected category to %s, actual %s", "+autostar", outline.Category)
	}

}

func getTestOPML() *OPML {
	// create opml programatically
	now := time.Date(2016, time.January, 14, 13, 58, 0, 0, time.Local).Truncate(time.Second)
	opml1 := &OPML{}
	opml1.Head = &Head{
		Title:       "Rakewire Subscriptions",
		DateCreated: now,
		OwnerName:   "karl@ostendorf.com",
	}
	opml1.Body = &Body{
		Outlines: Outlines{},
	}

	group1 := &Outline{
		Text:  "Group1",
		Title: "Group1",
	}
	group1.Outlines = append(group1.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g1text3",
		Title:   "g1title3",
		XMLURL:  "g1xmlurl3",
		HTMLURL: "g1htmlurl3",
	})
	group1.Outlines = append(group1.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g1text2",
		Title:   "g1title2",
		XMLURL:  "g1xmlurl2",
		HTMLURL: "g1htmlurl2",
	})
	group1.Outlines = append(group1.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g1text1",
		Title:   "g1title1",
		XMLURL:  "g1xmlurl1",
		HTMLURL: "g1htmlurl1",
	})

	group2 := &Outline{
		Text:  "Group2",
		Title: "Group2",
	}
	group2.Outlines = append(group2.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g2text3",
		Title:   "g2title3",
		XMLURL:  "g2xmlurl3",
		HTMLURL: "g2htmlurl3",
	})
	group2.Outlines = append(group2.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g2text2",
		Title:   "g2title2",
		XMLURL:  "g2xmlurl2",
		HTMLURL: "g2htmlurl2",
	})
	group2.Outlines = append(group2.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g2text1",
		Title:   "g2title1",
		XMLURL:  "g2xmlurl1",
		HTMLURL: "g2htmlurl1",
	})

	group3 := &Outline{
		Text:  "Group3",
		Title: "Group3",
	}

	group31 := &Outline{
		Text:  "Group31",
		Title: "Group31",
	}
	group31.Outlines = append(group31.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g31text3",
		Title:   "g31title3",
		XMLURL:  "g31xmlurl3",
		HTMLURL: "g31htmlurl3",
	})
	group31.Outlines = append(group31.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g31text2",
		Title:   "g31title2",
		XMLURL:  "g31xmlurl2",
		HTMLURL: "g31htmlurl2",
	})
	group31.Outlines = append(group31.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g31text1",
		Title:   "g31title1",
		XMLURL:  "g31xmlurl1",
		HTMLURL: "g31htmlurl1",
	})

	group32 := &Outline{
		Text:  "Group32",
		Title: "Group32",
	}
	group32.Outlines = append(group32.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g32text3",
		Title:   "g32title3",
		XMLURL:  "g32xmlurl3",
		HTMLURL: "g32htmlurl3",
	})
	group32.Outlines = append(group32.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g32text2",
		Title:   "g32title2",
		XMLURL:  "g32xmlurl2",
		HTMLURL: "g32htmlurl2",
	})
	group32.Outlines = append(group32.Outlines, &Outline{
		Created: &now,
		Type:    "rss",
		Text:    "g32text1",
		Title:   "g32title1",
		XMLURL:  "g32xmlurl1",
		HTMLURL: "g32htmlurl1",
	})

	group3.Outlines = append(group3.Outlines, group32)
	group3.Outlines = append(group3.Outlines, group31)

	opml1.Body.Outlines = append(opml1.Body.Outlines, group2)
	opml1.Body.Outlines = append(opml1.Body.Outlines, group1)
	opml1.Body.Outlines = append(opml1.Body.Outlines, group3)

	return opml1

}

func getTestDocument() io.Reader {
	document := `
	<?xml version="1.0" encoding="UTF-8"?>
	<opml>
		<head>
	    <title>Rakewire Subscriptions</title>
			<dateCreated>2016-01-14T13:58:00+01:00</dateCreated>
			<ownerName>karl@ostendorf.com</ownerName>
		</head>
		<body>
			<outline text="Group2" title="Group2">
				<outline type="rss" text="g2text3" title="g2title3" xmlUrl="g2xmlurl3" htmlUrl="g2htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g2text2" title="g2title2" xmlUrl="g2xmlurl2" htmlUrl="g2htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g2text1" title="g2title1" xmlUrl="g2xmlurl1" htmlUrl="g2htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
			</outline>
			<outline text="Group1" title="Group1">
				<outline type="rss" text="g1text3" title="g1title3" xmlUrl="g1xmlurl3" htmlUrl="g1htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g1text2" title="g1title2" xmlUrl="g1xmlurl2" htmlUrl="g1htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g1text1" title="g1title1" xmlUrl="g1xmlurl1" htmlUrl="g1htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
			</outline>
			<outline text="Group3" title="Group3">
				<outline text="Group32" title="Group32">
					<outline type="rss" text="g32text3" title="g32title3" xmlUrl="g32xmlurl3" htmlUrl="g32htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g32text2" title="g32title2" xmlUrl="g32xmlurl2" htmlUrl="g32htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g32text1" title="g32title1" xmlUrl="g32xmlurl1" htmlUrl="g32htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
				</outline>
				<outline text="Group31" title="Group31">
					<outline type="rss" text="g31text3" title="g31title3" xmlUrl="g31xmlurl3" htmlUrl="g31htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g31text2" title="g31title2" xmlUrl="g31xmlurl2" htmlUrl="g31htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g31text1" title="g31title1" xmlUrl="g31xmlurl1" htmlUrl="g31htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
				</outline>
			</outline>
		</body>
	</opml>`
	return strings.NewReader(document)
}

func getTestDocumentSorted() io.Reader {
	document := `
	<?xml version="1.0" encoding="UTF-8"?>
	<opml>
		<head>
	    <title>Rakewire Subscriptions</title>
			<dateCreated>2016-01-14T13:58:00+01:00</dateCreated>
			<ownerName>karl@ostendorf.com</ownerName>
		</head>
		<body>
			<outline text="Group1" title="Group1">
				<outline type="rss" text="g1text1" title="g1title1" xmlUrl="g1xmlurl1" htmlUrl="g1htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g1text2" title="g1title2" xmlUrl="g1xmlurl2" htmlUrl="g1htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g1text3" title="g1title3" xmlUrl="g1xmlurl3" htmlUrl="g1htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
			</outline>
			<outline text="Group2" title="Group2">
				<outline type="rss" text="g2text1" title="g2title1" xmlUrl="g2xmlurl1" htmlUrl="g2htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g2text2" title="g2title2" xmlUrl="g2xmlurl2" htmlUrl="g2htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
				<outline type="rss" text="g2text3" title="g2title3" xmlUrl="g2xmlurl3" htmlUrl="g2htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
			</outline>
			<outline text="Group3" title="Group3">
				<outline text="Group31" title="Group31">
					<outline type="rss" text="g31text1" title="g31title1" xmlUrl="g31xmlurl1" htmlUrl="g31htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g31text2" title="g31title2" xmlUrl="g31xmlurl2" htmlUrl="g31htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g31text3" title="g31title3" xmlUrl="g31xmlurl3" htmlUrl="g31htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
				</outline>
				<outline text="Group32" title="Group32">
					<outline type="rss" text="g32text1" title="g32title1" xmlUrl="g32xmlurl1" htmlUrl="g32htmlurl1" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g32text2" title="g32title2" xmlUrl="g32xmlurl2" htmlUrl="g32htmlurl2" created="2016-01-14T13:58:00+01:00"></outline>
					<outline type="rss" text="g32text3" title="g32title3" xmlUrl="g32xmlurl3" htmlUrl="g32htmlurl3" created="2016-01-14T13:58:00+01:00"></outline>
				</outline>
			</outline>
		</body>
	</opml>`
	return strings.NewReader(document)
}
