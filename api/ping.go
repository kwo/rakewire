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
	quitter := z.NewQuitter()
	done := ctx.Done()

sending:
	for {
		log.Debugf("sending ping")
		stream.Send(&pb.PingResponse{Time: time.Now().Unix()})
		select {
		case <-ticker.C:
		case <-quitter.C:
			break sending
		case <-done:
			break sending
		}
	}

	log.Debugf("exiting...")
	ticker.Stop()
	quitter.Stop()
	return nil

}
