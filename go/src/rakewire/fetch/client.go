package fetch

import (
	"errors"
	"net/http"
	"time"
)

// Internal Errors
var (
	ErrRedirected = errors.New("Fetch URL redirected")
)

func newInternalClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport:     &internalTransport{},
		CheckRedirect: noFlaggedRequestsRedirectPolicy,
		Timeout:       timeout,
	}
}

func noFlaggedRequestsRedirectPolicy(req *http.Request, via []*http.Request) error {
	if len(via) > 0 && via[0].Header.Get("X-Redirect") != "" {
		return ErrRedirected
	}
	return nil
}

type internalTransport struct {
	http.Transport
}

func (z *internalTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	rsp, err = z.Transport.RoundTrip(req)
	if err == nil {
		z.flagPermanentRedirects(req, rsp)
	}
	return
}

func (z *internalTransport) flagPermanentRedirects(req *http.Request, rsp *http.Response) {
	if rsp.StatusCode == http.StatusMovedPermanently {
		req.Header.Add("X-Redirect", "X")
	}
}
