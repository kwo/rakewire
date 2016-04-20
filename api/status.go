package api

import (
	"golang.org/x/net/context"
	"rakewire/api/pb"
)

// GetStatus implements the Status service.
func (z *API) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {

	_, errAuthorize := z.authorize(ctx)
	if errAuthorize != nil {
		return nil, errAuthorize
	}

	rsp := &pb.StatusResponse{
		Version:   z.version,
		BuildTime: z.buildTime,
		BuildHash: z.buildHash,
		AppStart:  z.appStart,
	}

	return rsp, nil

}
