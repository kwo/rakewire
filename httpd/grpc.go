package httpd

import (
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

// grpcHandler returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or handler otherwise. Copied from cockroachdb.
func grpcHandler(grpcServer *grpc.Server, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
