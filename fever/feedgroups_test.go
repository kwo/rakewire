package fever

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/kwo/rakewire/model"
	"strings"
	"testing"
)

func TestGroups(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)
	server := newServer(database)
	defer server.Close()

	user := getUser(t, database)
	var mSubscriptions model.Subscriptions
	err := database.Select(func(tx model.Transaction) error {
		mSubscriptions = model.S.GetForUser(tx, user.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot select subscriptions: %s", err.Error())
	}

	// make request
	target := server.URL + "/fever?api&groups"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if groups, err := response.GetObjectArray("groups"); err != nil {
		t.Fatalf("Error getting json groups: %s", err.Error())
	} else if len(groups) != 2 {
		t.Errorf("bad group count, expected %d, actual %d", 2, len(groups))
	} else {
		for i, group := range groups {
			if id, err := group.GetInt64("id"); err != nil {
				t.Errorf("Cannot retrieve group.id: %s", err.Error())
			} else if id <= 0 {
				t.Errorf("group.id mimatch, expected positive value, actual %d", id)
			}
			if name, err := group.GetString("title"); err != nil {
				t.Errorf("Cannot retrieve group.title: %s", err.Error())
			} else if name != fmt.Sprintf("Group%d", i) {
				t.Errorf("group.title mimatch, expected %s, actual %s", fmt.Sprintf("Group%d", i), name)
			}
		}
	}

	if feedGroups, err := response.GetObjectArray("feeds_groups"); err != nil {
		t.Fatalf("Error getting json feed_groups: %s", err.Error())
	} else if len(feedGroups) != 2 {
		t.Errorf("bad feed_group count, expected %d, actual %d", 2, len(feedGroups))
	} else {
		for i, feedGroup := range feedGroups {
			if id, err := feedGroup.GetInt64("group_id"); err != nil {
				t.Errorf("Cannot retrieve feed_group.group_id: %s", err.Error())
			} else if id <= 0 {
				t.Errorf("feed_group.group_id mimatch, expected positive value, actual %d", id)
			}
			if feedIDs, err := feedGroup.GetString("feed_ids"); err != nil {
				t.Errorf("Cannot retrieve feed_group.feed_ids: %s", err.Error())
			} else if len(feedIDs) == 0 {
				t.Error("feed_group.feed_ids is empty")
			} else {
				feedIDElements := strings.Split(feedIDs, ",")
				if len(feedIDElements) != 2 {
					t.Fatalf("bad FeedIDs size, expected %d elements, actual %d", 2, len(feedIDElements))
				} else {
					for j, feedIDElement := range feedIDElements {
						feedID := parseID(feedIDElement)
						if feedID == 0 {
							t.Errorf("Invalid FeedID: %s", err.Error())
						}
						feedIDStr := fmt.Sprintf("%010d", feedID)
						if feedIDStr != mSubscriptions[(j*2)+i].FeedID {
							t.Errorf("FeedID mismatch, expected %s, actual %s", mSubscriptions[(j*2)+i].FeedID, feedIDStr)
						}
					}
				}
			}
		}
	}

}

func TestFeeds(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)
	server := newServer(database)
	defer server.Close()

	user := getUser(t, database)
	var mSubscriptions model.Subscriptions
	err := database.Select(func(tx model.Transaction) error {
		mSubscriptions = model.S.GetForUser(tx, user.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("Cannot select subscriptions: %s", err.Error())
	}

	target := server.URL + "/fever?api&feeds"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if feeds, err := response.GetObjectArray("feeds"); err != nil {
		t.Fatalf("Error getting json feeds: %s", err.Error())
	} else if len(feeds) != 4 {
		t.Errorf("bad feed count, expected %d, actual %d", 4, len(feeds))
	} else {
		for i, feed := range feeds {
			if id, err := feed.GetInt64("id"); err != nil {
				t.Errorf("Cannot retrieve feed.id: %s", err.Error())
			} else if fmt.Sprintf("%010d", id) != mSubscriptions[i].FeedID {
				t.Errorf("feed.id mimatch, expected %s, actual %s", id, mSubscriptions[i].FeedID)
			}
			if title, err := feed.GetString("title"); err != nil {
				t.Errorf("Cannot retrieve feed.title: %s", err.Error())
			} else if title != mSubscriptions[i].Title {
				t.Errorf("feed.title mimatch, expected %s, actual %s", mSubscriptions[i].Title, title)
			}
		}
	}

	if feedGroups, err := response.GetObjectArray("feeds_groups"); err != nil {
		t.Fatalf("Error getting json feed_groups: %s", err.Error())
	} else if len(feedGroups) != 2 {
		t.Errorf("bad feed_group count, expected %d, actual %d", 2, len(feedGroups))
	} else {
		for i, feedGroup := range feedGroups {
			if id, err := feedGroup.GetInt64("group_id"); err != nil {
				t.Errorf("Cannot retrieve feed_group.group_id: %s", err.Error())
			} else if id <= 0 {
				t.Errorf("feed_group.group_id mimatch, expected positive value, actual %d", id)
			}
			if feedIDs, err := feedGroup.GetString("feed_ids"); err != nil {
				t.Errorf("Cannot retrieve feed_group.feed_ids: %s", err.Error())
			} else if len(feedIDs) == 0 {
				t.Error("feed_group.feed_ids is empty")
			} else {
				feedIDElements := strings.Split(feedIDs, ",")
				if len(feedIDElements) != 2 {
					t.Fatalf("bad FeedIDs size, expected %d elements, actual %d", 2, len(feedIDElements))
				} else {
					for j, feedIDElement := range feedIDElements {
						feedID := parseID(feedIDElement)
						if feedID == 0 {
							t.Errorf("Invalid FeedID: %s", err.Error())
						}
						feedIDStr := fmt.Sprintf("%010d", feedID)
						if feedIDStr != mSubscriptions[(j*2)+i].FeedID {
							t.Errorf("FeedID mismatch, expected %s, actual %s", mSubscriptions[(j*2)+i].FeedID, feedIDStr)
						}
					}
				}
			}
		}
	}

}
