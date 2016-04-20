package api

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	gateway "github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net/http"
	"rakewire/api/pb"
	"rakewire/model"
	"strings"
	"time"
)

const (
	authBasic = "Basic "
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

	pb.RegisterStatusServiceServer(grpcServer, z)

	ctx := context.Background()
	dcreds := credentials.NewTLS(tlsConfig)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := gateway.NewServeMux()
	if err := pb.RegisterStatusServiceHandlerFromEndpoint(ctx, gwmux, endpointConnect, dopts); err != nil {
		return nil, nil, err
	}

	router := http.NewServeMux()
	router.Handle("/", gwmux)

	return router, grpcServer, nil

}

func (z *API) authorize(ctx context.Context, roles ...string) (*model.User, error) {

	var authHeader string
	if md, ok := metadata.FromContext(ctx); ok {
		if auth, okAuth := md["authorization"]; okAuth {
			authHeader = auth[0]
		}
	}

	if authHeader == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	} else if strings.HasPrefix(authHeader, authBasic) {
		return z.authorizeBasic(authHeader, roles...)
	}

	return nil, grpc.Errorf(codes.Unauthenticated, fmt.Sprintf("%s: %s", codes.Unauthenticated.String(), "unknown authentication scheme"))

}

func (z *API) authorizeBasic(authHeader string, roles ...string) (*model.User, error) {

	errParse := grpc.Errorf(codes.Unauthenticated, fmt.Sprintf("%s: %s", codes.Unauthenticated.String(), "cannot parse basic authorization header"))
	errInternal := grpc.Errorf(codes.Internal, codes.Internal.String())
	errUnauthenticated := grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	errUnauthorized := grpc.Errorf(codes.PermissionDenied, codes.PermissionDenied.String())

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, errParse
	}

	var username, password string
	if data, err := base64.StdEncoding.DecodeString(fields[1]); err == nil {
		creds := strings.SplitN(string(data), ":", 2)
		if len(creds) != 2 {
			return nil, errParse
		}
		username = creds[0]
		password = creds[1]
	} else {
		return nil, errParse
	}

	// lookup user
	var user *model.User
	errDb := z.db.Select(func(tx model.Transaction) error {
		u := model.U.GetByUsername(tx, username)
		user = u
		return nil
	})
	if errDb != nil {
		return nil, errInternal
	}

	// check username
	if user == nil {
		return nil, errUnauthenticated
	}

	// check password
	if !user.MatchPassword(password) {
		return nil, errUnauthenticated
	}

	// check roles
	hasAllRoles := true
	for _, role := range roles {
		if !user.HasRole(role) {
			hasAllRoles = false
		}
		break
	}
	if !hasAllRoles {
		return nil, errUnauthorized
	}

	return user, nil

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
