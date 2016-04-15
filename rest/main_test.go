package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootNotFound(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	defer closeTestDatabase(t, database)

	apiRest := NewAPI("/api", database)
	server := httptest.NewServer(apiRest.Router())
	defer server.Close()

	u := server.URL + "/api"

	rsp, err := http.Get(u)
	if err != nil {
		t.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusNotFound {
		t.Fatalf("Bad error code, url: %s, expected %d, actual %d", u, http.StatusNotFound, rsp.StatusCode)
	}

}
