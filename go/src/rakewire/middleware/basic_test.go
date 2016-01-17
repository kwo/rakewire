package middleware

import (
	"encoding/base64"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rakewire/model"
	"testing"
)

func TestBasicAuthOK(t *testing.T) {

	t.Parallel()

	database, err := openDatabase()
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}
	server := getServer(database)
	defer server.Close()
	defer closeDatabase(database)

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Cannot construct request: %s", err.Error())
	}

	up := base64.StdEncoding.EncodeToString([]byte("karl:abcdefg"))
	req.Header.Set("Authorization", "Basic "+up)

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %s", err.Error())
	}

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("Bad status code, expected %d, actual %d", http.StatusOK, rsp.StatusCode)
	}

	expectedBody := "karl"
	if body, err := ioutil.ReadAll(rsp.Body); err != nil {
		t.Errorf("Error retrieving body: %s", err.Error())
	} else if string(body) != expectedBody {
		t.Errorf("Expected body %s, actual %s", expectedBody, body)
	}

}

func TestBasicAuthBadCredentials(t *testing.T) {

	t.Parallel()

	database, err := openDatabase()
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}
	server := getServer(database)
	defer server.Close()
	defer closeDatabase(database)

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Cannot construct request: %s", err.Error())
	}

	up := base64.StdEncoding.EncodeToString([]byte("karl:abcdef"))
	req.Header.Set("Authorization", "Basic "+up)

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %s", err.Error())
	}

	if rsp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Bad status code, expected %d, actual %d", http.StatusUnauthorized, rsp.StatusCode)
	}

}

func TestBasicAuthNoHeader(t *testing.T) {

	t.Parallel()

	// given

	database, err := openDatabase()
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}
	server := getServer(database)
	defer server.Close()
	defer closeDatabase(database)

	// when

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Cannot construct request: %s", err.Error())
	}

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %s", err.Error())
	}

	// then

	if rsp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Bad status code, expected %d, actual %d", http.StatusUnauthorized, rsp.StatusCode)
	}

}

func getServer(database model.Database) *httptest.Server {

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}

		if u := context.Get(req, "user"); u != nil {
			user := u.(*model.User)
			w.Write([]byte(user.Username))
			return
		}

		w.Write([]byte("KO\n"))

	})

	opts := &BasicAuthOptions{
		Realm:    "protected",
		Database: database,
	}

	return httptest.NewServer(Adapt(router, BasicAuth(opts)))

}
