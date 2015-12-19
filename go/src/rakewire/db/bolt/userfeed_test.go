package bolt

import (
	"fmt"
	m "rakewire/model"
	"testing"
)

func TestUserFeed(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	var users []*m.User
	for i := 0; i < 2; i++ {
		user := m.NewUser(fmt.Sprintf("User%d", i))
		user.SetPassword("abcdefg")
		if err := db.UserSave(user); err != nil {
			t.Fatalf("Error saving user: %s", err.Error())
		}
		users = append(users, user)
	}

	var feeds []*m.Feed
	for i := 0; i < 4; i++ {
		feed := m.NewFeed(fmt.Sprintf("http://localhost%d/", i))
		feed.Title = fmt.Sprintf("Feed%d", i)
		if err := db.SaveFeed(feed); err != nil {
			t.Fatalf("Error saving feed: %s", err.Error())
		}
		feeds = append(feeds, feed)
	}

	// save userfeeds
	for i := 0; i < 4; i++ {
		user := users[i%2]
		feed := feeds[i]
		userfeed := m.NewUserFeed(user.ID, feed.ID)
		if err := db.UserFeedSave(userfeed); err != nil {
			t.Fatalf("Error saving userfeed: %s", err.Error())
		}
	}

	for i := 0; i < 2; i++ {
		user := users[i]
		userfeeds, err := db.UserFeedGetAllByUser(user.ID)
		if err != nil {
			t.Fatalf("Error retrieving userfeeds: %s", err.Error())
		}
		if len(userfeeds) != 2 {
			t.Fatalf("bad userfeeds count, expected %d, actual %d", 2, len(userfeeds))
		}
		for j, userfeed := range userfeeds {
			feed := feeds[(j*2)+i]
			if userfeed.Feed == nil {
				t.Error("Feed not populated")
			}
			if userfeed.UserID != user.ID {
				t.Errorf("bad userID, expected %d, actual %d", user.ID, userfeed.UserID)
			}
			if userfeed.FeedID != feed.ID {
				t.Errorf("bad feedID, expected %d, actual %d", feed.ID, userfeed.FeedID)
			}
		}
	}

}
