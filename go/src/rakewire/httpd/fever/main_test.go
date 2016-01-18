package fever

import (
	"github.com/antonholmquist/jason"
	"rakewire/model"
	"strconv"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	server := newServer(database)
	defer server.Close()
	defer closeDatabase(t, database, databaseFile)

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

	target := server.URL + "/fever?api"
	data, err := makeRequest(user, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if v, err := response.GetInt64("api_version"); err != nil {
		t.Errorf("Error retrieving version from response: %s", err.Error())
	} else if v != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, v)
	}

	if v, err := response.GetInt64("auth"); err != nil {
		t.Errorf("Error retrieving authorized from response: %s", err.Error())
	} else if v != 1 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 1, v)
	}

	// check that last_refreshed_on_time is quoted
	if _, err := response.GetInt64("last_refreshed_on_time"); err == nil {
		t.Error("expected error getting last_refreshed_on_time as number, got none")
	}

	if v, err := response.GetString("last_refreshed_on_time"); err != nil {
		t.Errorf("last_refreshed_on_time not quoted as a string: %s", err.Error())
	} else if v == "" {
		t.Error("empty last_refreshed_on_time value")
	} else {
		if lastRefreshed, err := strconv.Atoi(v); err != nil {
			t.Errorf("Invalid value for last_refreshed_on_time: %s", err.Error())
		} else {
			uts := int(time.Now().Unix())
			if lastRefreshed < uts-1 || lastRefreshed > uts+1 {
				t.Errorf("last_refreshed_on_time mismatch, expected: %d, actual: %d", uts, lastRefreshed)
			}
		}
	}

}

func TestBadAuth(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)
	server := newServer(database)
	defer server.Close()

	target := server.URL + "/fever?api"
	data, err := makeRequest(nil, target)
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	if v, err := response.GetInt64("api_version"); err != nil {
		t.Errorf("Error retrieving version from response: %s", err.Error())
	} else if v != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, v)
	}

	if v, err := response.GetInt64("auth"); err != nil {
		t.Errorf("Error retrieving authorized from response: %s", err.Error())
	} else if v != 0 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 1, v)
	}

	if _, err := response.GetString("last_refreshed_on_time"); err == nil {
		t.Error("Error expected, got none, last_refreshed_on_time value is present on unauthorized access")
	}

}

func TestXmlUnsupported(t *testing.T) {

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

	target := server.URL + "/fever?api=xml"
	data, err := makeRequest(user, target)
	if err == nil {
		t.Fatal("Expected error but none")
	}
	t.Logf("raw response: %s", string(data))

}
