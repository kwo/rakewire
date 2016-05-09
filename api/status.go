package api

import (
	"github.com/kwo/rakewire/api/msg"
	"golang.org/x/net/context"
)

// GetStatus implements the Status service.
func (z *API) GetStatus(ctx context.Context, req *msg.StatusRequest) (*msg.StatusResponse, error) {

	rsp := &msg.StatusResponse{
		Version:   z.version,
		BuildTime: z.buildTime,
		BuildHash: z.buildHash,
		AppStart:  z.appStart,
	}

	return rsp, nil

}
