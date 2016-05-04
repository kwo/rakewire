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
	"sync"
	"time"
)

// API top level struct
type API struct {
	db          model.Database
	version     string
	buildTime   int64
	buildHash   string
	appStart    int64
	quitterLock sync.Mutex
	quitters    map[time.Time]chan bool
}

// New creates a new REST API instance
func New(database model.Database, versionString string, appStart time.Time) *API {

	version, buildTime, buildHash := parseVersionString(versionString)

	return &API{
		db:        database,
		version:   version,
		buildTime: buildTime,
		buildHash: buildHash,
		appStart:  appStart.Unix(),
		quitters:  make(map[time.Time]chan bool),
	}

}

// NewServer returns a new GRPC server
func (z *API) NewServer(tlsConfig *tls.Config) *grpc.Server {

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&tlsConfig.Certificates[0]))}
	rpc := grpc.NewServer(opts...)
	pb.RegisterPingServiceServer(rpc, z)
	pb.RegisterStatusServiceServer(rpc, z)
	pb.RegisterTokenServiceServer(rpc, z)

	return rpc

}

// NewHandler returns a new JSON gateway to a GRPC server
func (z *API) NewHandler(ctx context.Context, endpoint string, tlsConfig *tls.Config) (http.Handler, error) {

	dcreds := credentials.NewTLS(tlsConfig)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	router := gateway.NewServeMux()

	if err := pb.RegisterPingServiceHandlerFromEndpoint(ctx, router, endpoint, dopts); err != nil {
		return nil, err
	}
	if err := pb.RegisterStatusServiceHandlerFromEndpoint(ctx, router, endpoint, dopts); err != nil {
		return nil, err
	}
	if err := pb.RegisterTokenServiceHandlerFromEndpoint(ctx, router, endpoint, dopts); err != nil {
		return nil, err
	}

	return router, nil

}

// Shutdown gracefully shutdowns all connections to the API
func (z *API) Shutdown() {
	z.notifyQuitters()
	z.waitQuitters()
}

type fnDone func()

type quitter struct {
	C    <-chan bool
	Stop fnDone
}

func (z *API) newQuitter() *quitter {

	z.quitterLock.Lock()
	defer z.quitterLock.Unlock()

	id := time.Now()
	quitterChannel := make(chan bool)
	z.quitters[id] = quitterChannel

	done := func() {
		close(quitterChannel)
		z.quitterLock.Lock()
		defer z.quitterLock.Unlock()
		delete(z.quitters, id)
	}

	return &quitter{
		C:    quitterChannel,
		Stop: done,
	}

}

func (z *API) notifyQuitters() {
	z.quitterLock.Lock()
	defer z.quitterLock.Unlock()
	for _, quitterChannel := range z.quitters {
		quitterChannel <- true
	}
}

func (z *API) quitterCount() int {
	z.quitterLock.Lock()
	defer z.quitterLock.Unlock()
	return len(z.quitters)
}

func (z *API) waitQuitters() {
	for {
		if z.quitterCount() == 0 {
			break
		}
	}
}

func (z *API) authenticate(ctx context.Context, roles ...string) (*auth.User, error) {

	var authHeader string
	if md, ok := metadata.FromContext(ctx); ok {
		if header, okAuth := md["authorization"]; okAuth {
			authHeader = header[0]
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
