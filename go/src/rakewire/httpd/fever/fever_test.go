package fever

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"rakewire/db"
	"rakewire/db/bolt"
	"rakewire/logging"
	m "rakewire/model"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg := &logging.Configuration{Level: logging.LogWarn}
	cfg.Init()
	status := m.Run()
	os.Exit(status)
}

func TestAuthJson(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	// add test user
	user := m.NewUser("testuser@localhost")
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	apiFever := NewAPI("/fever", database)

	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	u := server.URL + "/fever?api"

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

	// check that last_refreshed_on_time is quoted
	lastRefreshOK := false
	dataFields := strings.Split(strings.Trim(dataString, " }{\r\n"), ",")
	if len(dataFields) != 3 {
		t.Fatal("raw json response is missing data")
	}
	for _, dataField := range dataFields {
		fields := strings.Split(dataField, ":")
		if len(fields) != 2 {
			t.Fatal("raw json response is missing data")
		}
		if fields[0] == "\"last_refreshed_on_time\"" {
			field01, err := strconv.Unquote(fields[1])
			if err != nil {
				t.Errorf("cannot unquote last_refreshed_on_time value (%s): %s", fields[1], err.Error())
				break
			}
			if len(field01) != len(fields[1])-2 {
				t.Errorf("last_refreshed_on_time is not quoted as a string: %s", fields[1])
			} else {
				lastRefreshOK = true
			}
		} // last refreshed
	}
	if !lastRefreshOK {
		t.Error("Cannot evaluate if last_refreshed_on_time is quoted as a string")
	}

	response := &Response{}
	if err := json.Unmarshal(data, response); err != nil {
		t.Fatalf("Invalid JSON response: %s\n", err.Error())
	}

	if response.Version != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, response.Version)
	}

	if response.Authorized != 1 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 1, response.Authorized)
	}

	uts := time.Now().Unix()
	if response.LastRefreshed < uts-1 || response.LastRefreshed > uts+1 {
		t.Errorf("LastRefreshed mismatch, expected: %d, actual: %d", uts, response.LastRefreshed)
	}

}

func TestAuthXml(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	// add test user
	user := m.NewUser("testuser@localhost")
	user.SetPassword("abcdefg")
	if err := database.UserSave(user); err != nil {
		t.Fatalf("Cannot save user: %s", err.Error())
	}

	apiFever := NewAPI("/fever", database)

	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	u := server.URL + "/fever?api=xml"

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
	if err := xml.Unmarshal(data, response); err != nil {
		t.Fatalf("Invalid XML response: %s\n", err.Error())
	}

	if response.Version != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, response.Version)
	}

	if response.Authorized != 1 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 1, response.Authorized)
	}

	uts := time.Now().Unix()
	if response.LastRefreshed < uts-1 || response.LastRefreshed > uts+1 {
		t.Errorf("LastRefreshed mismatch, expected: %d, actual: %d", uts, response.LastRefreshed)
	}

}

func TestBadAuthJson(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	apiFever := NewAPI("/fever", database)

	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	u := server.URL + "/fever?api"

	rsp, err := http.PostForm(u, nil)
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

	if ok := strings.Contains(dataString, "last_refreshed_on_time"); ok {
		t.Error("bad auth response cannot contain the last_refreshed_on_time field")
	}

	response := &Response{}
	if err := json.Unmarshal(data, response); err != nil {
		t.Fatalf("Invalid JSON response: %s\n", err.Error())
	}

	if response.Version != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, response.Version)
	}

	if response.Authorized != 0 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 0, response.Authorized)
	}

	if response.LastRefreshed != 0 {
		t.Errorf("LastRefreshed mismatch, expected: %d, actual: %d", 0, response.LastRefreshed)
	}

}

func TestBadAuthXml(t *testing.T) {

	t.Parallel()

	database, databaseFile := openDatabase(t)
	defer closeDatabase(t, database, databaseFile)

	apiFever := NewAPI("/fever", database)

	server := httptest.NewServer(apiFever.Router())
	defer server.Close()

	u := server.URL + "/fever?api=xml"

	rsp, err := http.PostForm(u, nil)
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

	if ok := strings.Contains(dataString, "last_refreshed_on_time"); ok {
		t.Error("bad auth response cannot contain the last_refreshed_on_time field")
	}

	response := &Response{}
	if err := xml.Unmarshal(data, response); err != nil {
		t.Fatalf("Invalid XML response: %s\n", err.Error())
	}

	if response.Version != 3 {
		t.Errorf("Version mismatch, expected: %d, actual: %d", 3, response.Version)
	}

	if response.Authorized != 0 {
		t.Errorf("Authorized mismatch, expected: %d, actual: %d", 0, response.Authorized)
	}

	if response.LastRefreshed != 0 {
		t.Errorf("LastRefreshed mismatch, expected: %d, actual: %d", 0, response.LastRefreshed)
	}

}

func openDatabase(t *testing.T) (*bolt.Service, string) {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Error creating tempfile: %s\n", err.Error())
	}
	testDatabaseFile := f.Name()
	f.Close()

	cfg := db.Configuration{
		Location: testDatabaseFile,
	}
	testDatabase := bolt.NewService(&cfg)
	err = testDatabase.Start()
	if err != nil {
		t.Fatalf("Cannot open database: %s\n", err.Error())
	}

	return testDatabase, testDatabaseFile

}

func closeDatabase(t *testing.T, database *bolt.Service, testDatabaseFile string) {

	database.Stop()
	if err := os.Remove(testDatabaseFile); err != nil {
		t.Errorf("Cannot delete temp database file: %s", err.Error())
	}

}
