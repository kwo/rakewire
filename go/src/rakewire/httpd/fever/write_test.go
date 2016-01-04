package fever

import (
	"github.com/antonholmquist/jason"
	"net/http/httptest"
	"testing"
)

func TestMark(t *testing.T) {

	t.SkipNow()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	user, err := database.UserGetByUsername(testUsername)
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

	// run server
	apiFever := NewAPI("/fever", database)
	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	// make request
	target := server.URL + "/fever?api"
	data, err := makeRequest(user, target, "mark", "item")
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}
	t.Logf("raw response: %s", string(data))

	response, err := jason.NewObjectFromBytes(data)
	if err != nil {
		t.Fatalf("Error parsing json response: %s", err.Error())
	}

	mark, err := response.GetString("mark")
	if err != nil {
		t.Fatalf("Cannot read mark attribute: %s", err.Error())
	} else if mark != "item" {
		t.Errorf("Bad mark attribute;: expected %s, actual %s", "item", mark)
	}

}
