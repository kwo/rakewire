package fever

import (
	"github.com/antonholmquist/jason"
	"rakewire/model"
	"testing"
)

func TestItemsAll(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)
	server := newServer(database)
	defer server.Close()

	var user *model.User
	err := database.Select(func(tx model.Transaction) error {
		u, err := model.UserByUsername(testUsername, tx)
		if err == nil && u != nil {
			user = u
		}
		return err
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var totalItems int64 = 40
	var expectedNumberItems = 40
	var expectedFirstID int64 = 1
	var expectedLastID int64 = 40

	// make request
	target := server.URL + "/fever?api&items"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if count, err := response.GetInt64("total_items"); err != nil {
		t.Fatalf("Cannot read total count: %s", err.Error())
	} else if count != totalItems {
		t.Errorf("Bad item count: expected %d, actual %d", totalItems, count)
	}

	items, err := response.GetObjectArray("items")
	if err != nil {
		t.Fatalf("Cannot read items: %s", err.Error())
	} else if len(items) != expectedNumberItems {
		t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(items))
	}

	if id, err := items[0].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedFirstID {
		t.Errorf("Bad first ID: expected %d, actual %d", expectedFirstID, id)
	}

	if id, err := items[len(items)-1].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedLastID {
		t.Errorf("Bad last ID: expected %d, actual %d", expectedLastID, id)
	}

}

func TestItemsNext(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)
	server := newServer(database)
	defer server.Close()

	var user *model.User
	err := database.Select(func(tx model.Transaction) error {
		u, err := model.UserByUsername(testUsername, tx)
		if err == nil && u != nil {
			user = u
		}
		return err
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var sinceID = "20"
	var totalItems int64 = 40
	var expectedNumberItems = 20
	var expectedFirstID int64 = 21
	var expectedLastID int64 = 40

	// make request
	target := server.URL + "/fever?api&items&since_id=" + sinceID
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if count, err := response.GetInt64("total_items"); err != nil {
		t.Fatalf("Cannot read total count: %s", err.Error())
	} else if count != totalItems {
		t.Errorf("Bad item count: expected %d, actual %d", totalItems, count)
	}

	items, err := response.GetObjectArray("items")
	if err != nil {
		t.Fatalf("Cannot read items: %s", err.Error())
	} else if len(items) != expectedNumberItems {
		t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(items))
	}

	if id, err := items[0].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedFirstID {
		t.Errorf("Bad first ID: expected %d, actual %d", expectedFirstID, id)
	}

	if id, err := items[len(items)-1].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedLastID {
		t.Errorf("Bad last ID: expected %d, actual %d", expectedLastID, id)
	}

}

func TestItemsPrev(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)
	server := newServer(database)
	defer server.Close()

	var user *model.User
	err := database.Select(func(tx model.Transaction) error {
		u, err := model.UserByUsername(testUsername, tx)
		if err == nil && u != nil {
			user = u
		}
		return err
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var maxID = "20"
	var totalItems int64 = 40
	var expectedNumberItems = 19
	var expectedFirstID int64 = 19
	var expectedLastID int64 = 1

	// make request
	target := server.URL + "/fever?api&items&max_id=" + maxID
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if count, err := response.GetInt64("total_items"); err != nil {
		t.Fatalf("Cannot read total count: %s", err.Error())
	} else if count != totalItems {
		t.Errorf("Bad item count: expected %d, actual %d", totalItems, count)
	}

	items, err := response.GetObjectArray("items")
	if err != nil {
		t.Fatalf("Cannot read items: %s", err.Error())
	} else if len(items) != expectedNumberItems {
		t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(items))
	}

	if id, err := items[0].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedFirstID {
		t.Errorf("Bad first ID: expected %d, actual %d", expectedFirstID, id)
	}

	if id, err := items[len(items)-1].GetInt64("id"); err != nil {
		t.Fatalf("Cannot read item ID: %s", err.Error())
	} else if id != expectedLastID {
		t.Errorf("Bad last ID: expected %d, actual %d", expectedLastID, id)
	}

}

func TestItemsByID(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)
	server := newServer(database)
	defer server.Close()

	var user *model.User
	err := database.Select(func(tx model.Transaction) error {
		u, err := model.UserByUsername(testUsername, tx)
		if err == nil && u != nil {
			user = u
		}
		return err
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	var withIDs = "15,37,8"
	var totalItems int64 = 40
	var expectedNumberItems = 3
	var expectedIDs = []int64{15, 37, 8}

	// make request
	target := server.URL + "/fever?api&items&with_ids=" + withIDs
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if count, err := response.GetInt64("total_items"); err != nil {
		t.Fatalf("Cannot read total count: %s", err.Error())
	} else if count != totalItems {
		t.Errorf("Bad item count: expected %d, actual %d", totalItems, count)
	}

	items, err := response.GetObjectArray("items")
	if err != nil {
		t.Fatalf("Cannot read items: %s", err.Error())
	} else if len(items) != expectedNumberItems {
		t.Errorf("Bad item count: expected %d, actual %d", expectedNumberItems, len(items))
	}

	for i, item := range items {
		if id, err := item.GetInt64("id"); err != nil {
			t.Fatalf("Cannot read item ID: %s", err.Error())
		} else if id != expectedIDs[i] {
			t.Errorf("Bad ID: expected %d, actual %d", expectedIDs[i], id)
		}

	}

}
