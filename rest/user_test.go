package rest

import (
	"log"
	"net/http"
	"net/http/httptest"
	"rakewire/model"
	"testing"
)

func TestUserGet(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	// add test user
	user := model.U.New("testuser")
	user.SetPassword("abcdefg")
	err := database.Update(func(tx model.Transaction) error {
		return model.U.Save(tx, user)
	})
	if err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	apiRest := NewAPI("/api", database)
	server := httptest.NewServer(apiRest.Router())
	defer server.Close()

	u := server.URL + "/api/users/" + user.Username

	rsp, err := http.Get(u)
	if err != nil {
		log.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusOK {
		t.Fatalf("Bad error code, url: %s, expected %d, actual %d", u, http.StatusOK, rsp.StatusCode)
	}

}
