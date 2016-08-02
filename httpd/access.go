package httpd

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kwo/rakewire/auth"
	"golang.org/x/net/context"
)

const (
	timestampKey = "httpT0"
)

// AccessTimerHandler returns a Handler which adds a timestamp to the context.
func AccessTimerHandler() Middleware {
	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ctx = context.WithValue(ctx, timestampKey, time.Now())
			next.ServeHTTPC(ctx, w, r)
		})
	}
}

// AccessLog logs http accesses to the console
func AccessLog() Middleware {
	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			rsp := &responseLogger{w: w}
			next.ServeHTTPC(ctx, rsp, r)
			writeLog(ctx, r, rsp.Status(), rsp.Size(), os.Stdout)
		})
	}
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func writeLog(ctx context.Context, r *http.Request, status, size int, w io.Writer) {

	username := "-"
	if user, ok := ctx.Value("user").(*auth.User); ok {
		username = user.Name
	} else if r.URL.User != nil {
		if name := r.URL.User.Username(); len(name) > 0 {
			username = name
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	uri := r.RequestURI
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}

	t1 := time.Now()
	t0 := t1
	if t, ok := ctx.Value(timestampKey).(time.Time); ok {
		t0 = t
	}

	ts := t0.Format("02/Jan/2006:15:04:05 -0700")
	elapsedTime := t1.Sub(t0)

	fmt.Fprintf(w, "%s - %s [%s] %s %s \"%s\" %d %d \"%s\" \"%s\" %s\n", host, username, ts, r.Method, uri, r.Proto, status, size, r.Referer(), r.UserAgent(), elapsedTime.String())

}
