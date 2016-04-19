package api

import (
	"crypto/tls"
	gateway "github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"rakewire/api/internal/pb"
	"rakewire/logger"
	"rakewire/model"
)

var (
	log = logger.New("api")
)

// NewAPI creates a new REST API instance
func NewAPI(database model.Database) *API {
	return &API{
		db: database,
	}
}

// API top level struct
type API struct {
	db model.Database
}

// Router returns the top-level router
func (z *API) Router(endpointConnect string, tlsConfig *tls.Config) (*http.ServeMux, *grpc.Server, error) {

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewServerTLSFromCert(&tlsConfig.Certificates[0]))}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEchoServiceServer(grpcServer, &echoer{})

	ctx := context.Background()
	dcreds := credentials.NewTLS(tlsConfig)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := gateway.NewServeMux()
	if err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, endpointConnect, dopts); err != nil {
		return nil, nil, err
	}

	router := http.NewServeMux()
	router.Handle("/", gwmux)

	return router, grpcServer, nil

}
