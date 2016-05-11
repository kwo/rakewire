package api

import (
	"github.com/kwo/rakewire/api/msg"
	"github.com/kwo/rakewire/auth"
	"golang.org/x/net/context"
)

// GetToken implements the Token service.
func (z *API) GetToken(ctx context.Context, req *msg.TokenRequest) (*msg.TokenResponse, error) {

	user := ctx.Value("user").(*auth.User)

	token, exp, errGenerate := auth.GenerateToken(user)
	if errGenerate != nil {
		return nil, errGenerate
	}

	rsp := &msg.TokenResponse{
		Username:   user.Name,
		Token:      string(token),
		Expiration: exp.Unix(),
	}

	return rsp, nil

}
