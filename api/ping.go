package api

import (
	"rakewire/api/pb"
	"rakewire/logger"
	"time"
)

// Ping implements the Ping service.
func (z *API) Ping(req *pb.PingRequest, stream pb.PingService_PingServer) error {

	// TODO: allow client to set interval

	ctx := stream.Context()
	_, errAuthorize := z.authenticate(ctx)
	if errAuthorize != nil {
		return errAuthorize
	}

	log := logger.New("ping")
	ticker := time.NewTicker(time.Second)
	done := ctx.Done()

sending:
	for {
		log.Debugf("sending ping")
		stream.Send(&pb.PingResponse{Time: time.Now().Unix()})
		select {
		case <-ticker.C:
		case <-done:
			break sending
		}
	}

	log.Debugf("exiting...")
	ticker.Stop()
	return nil

}