package native

import (
	"log"
	"net/http"
	"net/http/httptest"
	m "rakewire/model"
	"testing"
)

func TestUserGet(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	// add test user
	user := m.NewUser("testuser")
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	apiNative := NewAPI("/api", database)
	server := httptest.NewServer(apiNative.Router())
	defer server.Close()

	u := server.URL + "/api/users/" + user.Username

	rsp, err := http.Get(u)
	if err != nil {
		log.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusOK {
		t.Fatalf("Bad error code, url: %s, expected %d, actual %d", u, http.StatusOK, rsp.StatusCode)
	}

}
