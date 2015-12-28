package native

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootNotFound(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	apiNative := NewAPI("/api", database)
	server := httptest.NewServer(apiNative.Router())
	defer server.Close()

	u := server.URL + "/api"

	rsp, err := http.Get(u)
	if err != nil {
		log.Fatalf("Cannot perform request to %s: %s", u, err.Error())
	} else if rsp.StatusCode != http.StatusNotFound {
		t.Fatalf("Bad error code, url: %s, expected %d, actual %d", u, http.StatusNotFound, rsp.StatusCode)
	}

}
