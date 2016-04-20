package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"rakewire/api/pb"
)

type echoer struct{}

// Echo implements the Echo service.
func (z *echoer) Echo(ctx context.Context, msg *pb.EchoMessage) (*pb.EchoMessage, error) {
	if md, ok := metadata.FromContext(ctx); ok {
		if auth, okAuth := md["authorization"]; okAuth {
			log.Debugf("echo metadata auth: %s", auth[0])
		}
	}
	log.Debugf("echo request: %s", msg.Value)
	return msg, nil
}
