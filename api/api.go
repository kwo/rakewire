package api

import (
	"crypto/tls"
	"fmt"
	gateway "github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net/http"
	"rakewire/api/pb"
	"rakewire/auth"
	"rakewire/model"
	"strings"
	"time"
)

// NewAPI creates a new REST API instance
func NewAPI(database model.Database, versionString string, appStart time.Time) *API {

	version, buildTime, buildHash := parseVersionString(versionString)

	return &API{
		db:        database,
		version:   version,
		buildTime: buildTime,
		buildHash: buildHash,
		appStart:  appStart.Unix(),
	}

}

// API top level struct
type API struct {
	db        model.Database
	version   string
	buildTime int64
	buildHash string
	appStart  int64
}

// Router returns the top-level router
func (z *API) Router(endpointConnect string, tlsConfig *tls.Config) (http.Handler, *grpc.Server, error) {

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&tlsConfig.Certificates[0]))}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterPingServiceServer(grpcServer, z)
	pb.RegisterStatusServiceServer(grpcServer, z)

	ctx := context.Background()
	dcreds := credentials.NewTLS(tlsConfig)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := gateway.NewServeMux()

	if err := pb.RegisterPingServiceHandlerFromEndpoint(ctx, gwmux, endpointConnect, dopts); err != nil {
		return nil, nil, err
	}

	if err := pb.RegisterStatusServiceHandlerFromEndpoint(ctx, gwmux, endpointConnect, dopts); err != nil {
		return nil, nil, err
	}

	router := http.NewServeMux()
	router.Handle("/", gwmux)

	return router, grpcServer, nil

}

func (z *API) authenticate(ctx context.Context, roles ...string) (*model.User, error) {

	var authHeader string
	if md, ok := metadata.FromContext(ctx); ok {
		if auth, okAuth := md["authorization"]; okAuth {
			authHeader = auth[0]
		}
	}

	if user, errAuth := auth.Authenticate(z.db, authHeader, roles...); errAuth == nil {
		return user, nil
	} else if errAuth == auth.ErrUnauthenticated {
		return nil, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	} else if errAuth == auth.ErrUnauthorized {
		return nil, grpc.Errorf(codes.PermissionDenied, codes.PermissionDenied.String())
	} else if errAuth == auth.ErrBadHeader {
		return nil, grpc.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
	} else {
		return nil, grpc.Errorf(codes.Internal, fmt.Sprintf("%s: %s", codes.Internal.String(), errAuth.Error()))
	}

}

func parseVersionString(versionString string) (string, int64, string) {

	// parse version string
	fields := strings.Fields(versionString)
	if len(fields) == 3 {

		version := fields[0]
		buildHash := fields[2]

		var buildTime int64
		if bt, err := time.Parse(time.RFC3339, fields[1]); err == nil {
			buildTime = bt.Unix()
		}

		return version, buildTime, buildHash

	}

	return "", 0, ""

}
