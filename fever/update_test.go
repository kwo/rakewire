package fever

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/kwo/rakewire/model"
	"testing"
)

func TestMark(t *testing.T) {

	t.SkipNow()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)
	server := newServer(database)
	defer server.Close()

	var user *model.User
	err := database.Select(func(tx model.Transaction) error {
		u := model.U.GetByUsername(tx, testUsername)
		if u != nil {
			user = u
		}
		return fmt.Errorf("User unknown: %s", testUsername)
	})
	if err != nil {
		t.Fatalf("Cannot get user: %s", err.Error())
	}

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
