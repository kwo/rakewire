package fever

import (
	"github.com/antonholmquist/jason"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestUnreadIDs(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	user, err := database.UserGetByUsername(testUsername)
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var expectedNumberItems = 24
	// var expectedFirstID uint64 = 1
	// var expectedLastID uint64 = 40

	// run server
	apiFever := NewAPI("/fever", database)
	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	// make request
	target := server.URL + "/fever?api&unread_item_ids"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	itemIDStr, err := response.GetString("unread_item_ids")
	if err != nil {
		t.Fatalf("Cannot read unread item ids: %s", err.Error())
	} else if itemIDStr == "" {
		t.Fatal("Blank item IDs")
	} else {

		itemIDStrArray := strings.Split(itemIDStr, ",")
		if len(itemIDStrArray) != expectedNumberItems {
			t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(itemIDStrArray))
		}

		itemIDs := []uint64{}
		for _, idStr := range itemIDStrArray {
			if id, err := strconv.ParseUint(idStr, 10, 64); err != nil {
				t.Errorf("Cannot parse ID (%s): %s", idStr, err.Error())
			} else {
				itemIDs = append(itemIDs, id)
			}
		}

		// for _, itemID := range itemIDs {
		// 	t.Logf("itemID: %d", itemID)
		// }

		if len(itemIDs) != expectedNumberItems {
			t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(itemIDs))
		}

		// if itemIDs[0] != expectedFirstID {
		// 	t.Errorf("Bad first ID: expected %d, actual %d", expectedFirstID, itemIDs[0])
		// }
		//
		// if itemIDs[len(itemIDs)-1] != expectedLastID {
		// 	t.Errorf("Bad last ID: expected %d, actual %d", expectedLastID, itemIDs[len(itemIDs)-1])
		// }

	}

}

func TestSavedIDs(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	user, err := database.UserGetByUsername(testUsername)
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var expectedNumberItems = 8
	// var expectedFirstID uint64 = 1
	// var expectedLastID uint64 = 40

	// run server
	apiFever := NewAPI("/fever", database)
	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	// make request
	target := server.URL + "/fever?api&saved_item_ids"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	itemIDStr, err := response.GetString("saved_item_ids")
	if err != nil {
		t.Fatalf("Cannot read saved item ids: %s", err.Error())
	} else if itemIDStr == "" {
		t.Fatal("Blank item IDs")
	} else {

		itemIDStrArray := strings.Split(itemIDStr, ",")
		if len(itemIDStrArray) != expectedNumberItems {
			t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(itemIDStrArray))
		}

		itemIDs := []uint64{}
		for _, idStr := range itemIDStrArray {
			if id, err := strconv.ParseUint(idStr, 10, 64); err != nil {
				t.Errorf("Cannot parse ID (%s): %s", idStr, err.Error())
			} else {
				itemIDs = append(itemIDs, id)
			}
		}

		if len(itemIDs) != expectedNumberItems {
			t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(itemIDs))
		}

		// if itemIDs[0] != expectedFirstID {
		// 	t.Errorf("Bad first ID: expected %d, actual %d", expectedFirstID, itemIDs[0])
		// }
		//
		// if itemIDs[len(itemIDs)-1] != expectedLastID {
		// 	t.Errorf("Bad last ID: expected %d, actual %d", expectedLastID, itemIDs[len(itemIDs)-1])
		// }

	}

}
