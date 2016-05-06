package api

import (
	"golang.org/x/net/context"
	"time"
)

// StatusRequest defines the status request
type StatusRequest struct{}

// StatusResponse defines the status response
type StatusResponse struct {
	Version   string    `json:"version"`
	BuildTime time.Time `json:"buildTime"`
	BuildHash string    `json:"buildHash"`
	AppStart  time.Time `json:"appStart"`
}

// GetStatus implements the Status service.
func (z *API) GetStatus(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {

	rsp := &StatusResponse{
		Version:   z.version,
		BuildTime: z.buildTime,
		BuildHash: z.buildHash,
		AppStart:  z.appStart,
	}

	return rsp, nil

}
