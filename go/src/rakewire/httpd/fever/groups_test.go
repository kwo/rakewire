package fever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	m "rakewire/model"
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
	for i := 0; i < 2; i++ {
		g := m.NewGroup(user.ID, fmt.Sprintf("Group%d", i))
		if err := database.GroupSave(g); err != nil {
			t.Fatalf("Cannot add group: %s", err.Error())
		}
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

	dataString := string(data)
	t.Logf("raw response: %s", dataString)

	response := &Response{}
	if err := json.Unmarshal(data, response); err != nil {
		t.Fatalf("Invalid JSON response: %s\n", err.Error())
	}

	if response.Groups == nil {
		t.Fatal("No groups")
	}

	if len(response.Groups) != 2 {
		t.Errorf("bad group count, expected %d, actual %d", 2, len(response.Groups))
	}

}
