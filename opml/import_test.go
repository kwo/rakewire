package opml

import (
	"io"
	"rakewire/model"
	"strings"
	"testing"
	"time"
)

func TestImportAutoStar(t *testing.T) {

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)
	user := addUser(t, db)

	importData(t, db, user, getImport1())
	verifyAutoStar(t, db, user, "g1title2", true)

	importData(t, db, user, getImport2())
	verifyAutoStar(t, db, user, "g1title2", false)

}

func TestImportCreatedDate(t *testing.T) {

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)
	user := addUser(t, db)

	now := time.Now().Truncate(time.Second)
	importData(t, db, user, getImport1())
	verifyCreatedDate(t, db, user, "g1title1", now)

	importData(t, db, user, getImport2())
	verifyCreatedDate(t, db, user, "g1title1", now)

}

func TestImportFeeds(t *testing.T) {

	feeds := []string{
		"g1xmlurl1", "g1xmlurl2", "g1xmlurl3",
		"g2xmlurl1", "g2xmlurl2", "g2xmlurl3",
		"g3xmlurl1", "g3xmlurl2", "g3xmlurl3",
		"g4xmlurl1", "g4xmlurl2", "g4xmlurl3",
	}

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)
	user := addUser(t, db)

	importData(t, db, user, getImport1())
	verifyFeeds(t, db, user, feeds)

	importData(t, db, user, getImport2())
	verifyFeeds(t, db, user, feeds)

}

func TestImportSubscriptions(t *testing.T) {

	before := map[string][]string{
		"Group1":        {"g1title1", "g1title2", "g1title3"},
		"Group2":        {"g2title1", "g2title2", "g2title3"},
		"GroupX/Group3": {"g3title1", "g3title2", "g3title3"},
		"GroupX/Group4": {"g4title1", "g4title2", "g4title3"},
	}

	after := map[string][]string{
		"Group1":        {"g1title1", "g1title2", "g2title1"},
		"Group2":        {"g1title3", "g2title2", "g2title3"},
		"GroupX/Group3": {},
		"GroupX/Group4": {},
	}

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)
	user := addUser(t, db)

	importData(t, db, user, getImport1())
	verifyGroups(t, db, user, extractMapKeys(before))
	verifySubscriptions(t, db, user, "", extractMapValues(before))
	for key, values := range before {
		verifySubscriptions(t, db, user, key, values)
	}

	importData(t, db, user, getImport2())
	verifyGroups(t, db, user, extractMapKeys(after))
	verifySubscriptions(t, db, user, "", extractMapValues(after))
	for key, values := range after {
		verifySubscriptions(t, db, user, key, values)
	}

}

func verifyAutoStar(t *testing.T, db model.Database, user *model.User, title string, flag bool) {

	err := db.Select(func(tx model.Transaction) error {

		subscription := model.S.GetForUser(tx, user.ID).ByTitle()[title]

		if subscription.AutoStar != flag {
			t.Errorf("Bad autostar flag: %t, expected %t", subscription.AutoStar, flag)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error verifying feeds: %s", err.Error())
	}

}

func verifyCreatedDate(t *testing.T, db model.Database, user *model.User, title string, importTime time.Time) {

	err := db.Select(func(tx model.Transaction) error {

		subscription := model.S.GetForUser(tx, user.ID).ByTitle()[title]

		if subscription.Added.Before(importTime) {
			t.Errorf("Import should ignore created dates on subscriptions: %s > %s", importTime.Format(time.RFC3339), subscription.Added.Format(time.RFC3339))
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error verifying feeds: %s", err.Error())
	}

}

func verifyFeeds(t *testing.T, db model.Database, user *model.User, expectedFeeds []string) {

	err := db.Select(func(tx model.Transaction) error {

		feeds := model.F.GetNext(tx, time.Now().Add(24*time.Hour))

		expectedFeedCount := len(expectedFeeds)
		feedCount := len(feeds)
		if feedCount != expectedFeedCount {
			t.Errorf("Bad feed count: %d, expected %d", feedCount, expectedFeedCount)
		}

		byURL := feeds.ByURL()
		for _, expectedURL := range expectedFeeds {
			if _, ok := byURL[expectedURL]; !ok {
				t.Errorf("Missing feed: %s", expectedURL)
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error verifying feeds: %s", err.Error())
	}

}

func verifyGroups(t *testing.T, db model.Database, user *model.User, expectedGroups []string) {

	err := db.Select(func(tx model.Transaction) error {

		groups := model.G.GetForUser(tx, user.ID)
		expectedGroupCount := len(expectedGroups)
		groupCount := len(groups)
		if groupCount != expectedGroupCount {
			t.Errorf("Bad group count: %d, expected %d", groupCount, expectedGroupCount)
		}

		groupsByName := groups.ByName()
		for _, expectedName := range expectedGroups {
			if _, ok := groupsByName[expectedName]; !ok {
				t.Errorf("Missing group: %s", expectedName)
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error verifying groups: %s", err.Error())
	}

}

func verifySubscriptions(t *testing.T, db model.Database, user *model.User, groupName string, expectedSubscriptions []string) {

	err := db.Select(func(tx model.Transaction) error {

		subscriptions := model.S.GetForUser(tx, user.ID)
		if groupName != "" {
			group := model.G.GetForUser(tx, user.ID).ByName()[groupName]
			subscriptions = subscriptions.WithGroup(group.ID)
		}

		expectedSubscriptionCount := len(expectedSubscriptions)
		subscriptionCount := len(subscriptions)
		if subscriptionCount != expectedSubscriptionCount {
			t.Errorf("Bad subscription count: %d, expected %d", subscriptionCount, expectedSubscriptionCount)
		}

		subscriptionsByTitle := subscriptions.ByTitle()
		for _, expectedTitle := range expectedSubscriptions {
			if _, ok := subscriptionsByTitle[expectedTitle]; !ok {
				t.Errorf("Missing subscription: %s", expectedTitle)
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error verifying subscriptions: %s", err.Error())
	}

}

func importData(t *testing.T, db model.Database, user *model.User, reader io.Reader) {

	err := db.Update(func(tx model.Transaction) error {

		// parse opml
		o, err := Parse(reader)
		if err != nil {
			return err
		}

		// import
		if err := Import(tx, user.ID, o); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error importing OPML: %s", err.Error())
	}

}

func extractMapKeys(m map[string][]string) []string {
	var result []string
	for key := range m {
		result = append(result, key)
	}
	return result
}

func extractMapValues(m map[string][]string) []string {
	var result []string
	for _, a := range m {
		for _, value := range a {
			result = append(result, value)
		}
	}
	return result
}

func getImport1() io.Reader {
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
				<outline type="rss" title="g2title3" xmlUrl="g2xmlurl3" htmlUrl="g2htmlurl3"/>
				<outline type="rss" title="g2title2" xmlUrl="g2xmlurl2" htmlUrl="g2htmlurl2"/>
				<outline type="rss" title="g2title1" xmlUrl="g2xmlurl1" htmlUrl="g2htmlurl1"/>
			</outline>
			<outline text="Group1" title="Group1">
				<outline type="rss" title="g1title3" xmlUrl="g1xmlurl3" htmlUrl="g1htmlurl3"/>
				<outline type="rss" title="g1title2" xmlUrl="g1xmlurl2" htmlUrl="g1htmlurl2" category="+autostar" />
				<outline type="rss" title="g1title1" xmlUrl="g1xmlurl1" htmlUrl="g1htmlurl1" created="2016-01-01T00:00:00Z"/>
			</outline>
			<outline text="GroupX" title="GroupX">
				<outline text="Group3" title="Group3">
					<outline type="rss" title="g3title3" xmlUrl="g3xmlurl3" htmlUrl="g3htmlurl3"/>
					<outline type="rss" title="g3title2" xmlUrl="g3xmlurl2" htmlUrl="g3htmlurl2"/>
					<outline type="rss" title="g3title1" xmlUrl="g3xmlurl1" htmlUrl="g3htmlurl1"/>
				</outline>
				<outline text="Group4" title="Group4">
					<outline type="rss" text="g4title3" xmlUrl="g4xmlurl3" htmlUrl="g4htmlurl3"/>
					<outline type="rss" text="g4title2" xmlUrl="g4xmlurl2" htmlUrl="g4htmlurl2"/>
					<outline type="rss" text="g4title1" xmlUrl="g4xmlurl1" htmlUrl="g4htmlurl1"/>
				</outline>
			</outline>
		</body>
	</opml>`
	return strings.NewReader(document)
}

func getImport2() io.Reader {
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
				<outline type="rss" title="g2title3" xmlUrl="g2xmlurl3" htmlUrl="g2htmlurl3"/>
				<outline type="rss" title="g2title2" xmlUrl="g2xmlurl2" htmlUrl="g2htmlurl2"/>
				<outline type="rss" title="g1title3" xmlUrl="g1xmlurl3" htmlUrl="g1htmlurl3"/>
			</outline>
			<outline text="Group1" title="Group1">
				<outline type="rss" title="g2title1" xmlUrl="g2xmlurl1" htmlUrl="g2htmlurl1"/>
				<outline type="rss" title="g1title2" xmlUrl="g1xmlurl2" htmlUrl="g1htmlurl2"/>
				<outline type="rss" title="g1title1" xmlUrl="g1xmlurl1" htmlUrl="g1htmlurl1" created="2016-01-01T00:00:00Z"/>
			</outline>
		</body>
	</opml>`
	return strings.NewReader(document)
}
