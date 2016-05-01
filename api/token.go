package api

import (
	"golang.org/x/net/context"
	"rakewire/api/pb"
	"rakewire/auth"
)

// GetToken implements the Token service.
func (z *API) GetToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {

	user, errAuthorize := z.authenticate(ctx)
	if errAuthorize != nil {
		return nil, errAuthorize
	}

	token, exp, errGenerate := auth.GenerateToken(user)
	if errGenerate != nil {
		return nil, errGenerate
	}

	rsp := &pb.TokenResponse{
		Token:      string(token),
		Expiration: exp.Unix(),
	}

	return rsp, nil

}
