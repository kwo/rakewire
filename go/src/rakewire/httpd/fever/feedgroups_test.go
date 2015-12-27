package fever

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	m "rakewire/model"
	"strconv"
	"strings"
	"testing"
)

func TestGroups(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	// add test user
	user := m.NewUser("testuser@localhost")
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	// add test groups
	mGroups := []*m.Group{}
	for i := 0; i < 2; i++ {
		g := m.NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := database.GroupSave(g); err != nil {
			t.Fatalf("Cannot add group: %s", err.Error())
		}
	}

	// add test feeds
	mUserFeeds := []*m.UserFeed{}
	for i := 0; i < 4; i++ {
		f := m.NewFeed(fmt.Sprintf("http://localhost%d", i))
		if _, err := database.FeedSave(f); err != nil {
			t.Fatalf("Cannot add feed: %s", err.Error())
		}
		uf := m.NewUserFeed(user.ID, f.ID)
		uf.GroupIDs = append(uf.GroupIDs, mGroups[i%2].ID)
		database.UserFeedSave(uf)
		mUserFeeds = append(mUserFeeds, uf)
	}

	// run server
	apiFever := NewAPI("/fever", database)
	server := httptest.NewServer(apiFever.Router())
	defer server.Close()
	u := server.URL + "/fever?api&groups"

	values := url.Values{}
	values.Set(AuthParam, user.FeverHash)
	rsp, err := http.PostForm(u, values)
	if err != nil {
		log.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusOK {
		t.Fatalf("Bad error code, expected %d, actual %d", http.StatusOK, rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
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

	if feedGroups, err := response.GetObjectArray("feed_groups"); err != nil {
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
			} else if feedIDs == "" {
				t.Error("feed_group.feed_ids is empty")
			} else {
				feedIDElements := strings.Split(feedIDs, ",")
				if len(feedIDElements) != 2 {
					t.Fatalf("bad FeedIDs size, expected %d elements, actual %d", 2, len(feedIDElements))
				} else {
					for j, feedIDElement := range feedIDElements {
						feedID, err := strconv.Atoi(feedIDElement)
						if err != nil {
							t.Errorf("Invalid FeedID: %s", err.Error())
						}
						if uint64(feedID) != mUserFeeds[(j*2)+i].ID {
							t.Errorf("FeedID mismatch, expected %d, actual %d", mUserFeeds[(j*2)+i].ID, feedID)
						}
					}
				}
			}
		}
	}

}

func TestFeeds(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	// add test user
	user := m.NewUser("testuser@localhost")
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	// add test groups
	mGroups := []*m.Group{}
	for i := 0; i < 2; i++ {
		g := m.NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		mGroups = append(mGroups, g)
		if err := database.GroupSave(g); err != nil {
			t.Fatalf("Cannot add group: %s", err.Error())
		}
	}

	// add test feeds
	mUserFeeds := []*m.UserFeed{}
	for i := 0; i < 4; i++ {
		f := m.NewFeed(fmt.Sprintf("http://localhost%d", i))
		if _, err := database.FeedSave(f); err != nil {
			t.Fatalf("Cannot add feed: %s", err.Error())
		}
		uf := m.NewUserFeed(user.ID, f.ID)
		uf.Title = fmt.Sprintf("UserFeed%d", i)
		uf.GroupIDs = append(uf.GroupIDs, mGroups[i%2].ID)
		database.UserFeedSave(uf)
		mUserFeeds = append(mUserFeeds, uf)
	}

	// run server
	apiFever := NewAPI("/fever", database)
	server := httptest.NewServer(apiFever.Router())
	defer server.Close()
	u := server.URL + "/fever?api&feeds"

	values := url.Values{}
	values.Set(AuthParam, user.FeverHash)
	rsp, err := http.PostForm(u, values)
	if err != nil {
		log.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusOK {
		t.Fatalf("Bad error code, expected %d, actual %d", http.StatusOK, rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
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
			} else if uint64(id) != mUserFeeds[i].ID {
				t.Errorf("feed.id mimatch, expected %d, actual %d", id, mUserFeeds[i].ID)
			}
			if title, err := feed.GetString("title"); err != nil {
				t.Errorf("Cannot retrieve feed.title: %s", err.Error())
			} else if title != mUserFeeds[i].Title {
				t.Errorf("feed.title mimatch, expected %s, actual %s", mUserFeeds[i].Title, title)
			}
		}
	}

	if feedGroups, err := response.GetObjectArray("feed_groups"); err != nil {
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
			} else if feedIDs == "" {
				t.Error("feed_group.feed_ids is empty")
			} else {
				feedIDElements := strings.Split(feedIDs, ",")
				if len(feedIDElements) != 2 {
					t.Fatalf("bad FeedIDs size, expected %d elements, actual %d", 2, len(feedIDElements))
				} else {
					for j, feedIDElement := range feedIDElements {
						feedID, err := strconv.Atoi(feedIDElement)
						if err != nil {
							t.Errorf("Invalid FeedID: %s", err.Error())
						}
						if uint64(feedID) != mUserFeeds[(j*2)+i].ID {
							t.Errorf("FeedID mismatch, expected %d, actual %d", mUserFeeds[(j*2)+i].ID, feedID)
						}
					}
				}
			}
		}
	}

}
