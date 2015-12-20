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

}
