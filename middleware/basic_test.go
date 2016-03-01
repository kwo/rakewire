package middleware

import (
	"encoding/base64"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"rakewire/model"
	"testing"
)

func TestBasicAuthOK(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t)
	server := getServer(database)
	defer server.Close()
	defer closeTestDatabase(t, database)

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

	database := openTestDatabase(t)
	server := getServer(database)
	defer server.Close()
	defer closeTestDatabase(t, database)

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

	database := openTestDatabase(t)
	server := getServer(database)
	defer server.Close()
	defer closeTestDatabase(t, database)

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

func openTestDatabase(t *testing.T) model.Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	boltDB, err := model.OpenDatabase(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	err = boltDB.Update(func(tx model.Transaction) error {
		return populateDatabase(tx)
	})
	if err != nil {
		t.Fatalf("Cannot populate database: %s", err.Error())
	}

	return boltDB

}

func closeTestDatabase(t *testing.T, d model.Database) {

	location := d.Location()

	if err := model.CloseDatabase(d); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func populateDatabase(tx model.Transaction) error {

	// add test user
	user := model.NewUser("karl")
	user.SetPassword("abcdefg")
	return user.Save(tx)

}
