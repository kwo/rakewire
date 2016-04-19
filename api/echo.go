package api

import (
	"golang.org/x/net/context"
	"rakewire/api/internal/pb"
)

type echoer struct{}

// Echo implements the Echo service.
func (z *echoer) Echo(ctx context.Context, msg *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Debugf("rpc request Echo(%q)", msg.Value)
	return msg, nil
}
