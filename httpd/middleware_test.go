package httpd

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {

	handlerC := Chain(endstation(), fakeware("2"), fakeware("1"))
	handler := Adapt(context.WithValue(context.Background(), "a", []string{}), handlerC)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Cannot create request: %s", err.Error())
	}

	rsp, errRsp := client.Do(req)
	if errRsp != nil {
		t.Errorf("Error getting response: %s", errRsp.Error())
	}

	body, errBody := ioutil.ReadAll(rsp.Body)
	if errBody != nil {
		t.Errorf("Error reading response: %s", errBody.Error())
	}
	defer rsp.Body.Close()
	bodyStr := string(body)

	expectedOrder := "1 2"
	if bodyStr != expectedOrder {
		t.Errorf("Unexpected ordering of middleware: %s, expected: %s", bodyStr, expectedOrder)
	}

	t.Logf("Response: %s", bodyStr)

}

func endstation() HandlerC {
	return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if a, ok := ctx.Value("a").([]string); ok {
			w.Write([]byte(strings.Join(a, " ")))
		}
	})
}

func fakeware(id string) Middleware {
	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if a, ok := ctx.Value("a").([]string); ok {
				a = append(a, id)
				ctx = context.WithValue(ctx, "a", a)
			}
			next.ServeHTTPC(ctx, w, r)
		})
	}
}
