package httpd

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func newServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Host: %s, Path: %s, URL: %s", r.Host, r.URL.Path, r.URL.String())
		t.Logf("URL: %#v", r.URL)
		w.Write([]byte("OK\n"))
	}))
}

func TestRequest(t *testing.T) {

	t.SkipNow()

	server := newServer(t)
	defer server.Close()

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.URL+"/abc/123.html?this=that&what=what", nil)
	if err != nil {
		t.Fatalf("Cannot create request: %s", err.Error())
	}
	req.Header.Add("Host", "serverville")

	client.Do(req)

}

func TestHTML5(t *testing.T) {
	r := regexp.MustCompile("^/[a-z0-9/-]+$")
	if !r.MatchString("/index") {
		t.Error("Expected /index to match")
	}
	if r.MatchString("/index.html") {
		t.Error("Expected /index.html not to match")
	}
}
