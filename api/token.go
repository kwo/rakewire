package api

import (
	"golang.org/x/net/context"
	"rakewire/auth"
	"time"
)

// TokenRequest defines the token request
type TokenRequest struct{}

// TokenResponse defines the token response
type TokenResponse struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

// GetToken implements the Token service.
func (z *API) GetToken(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {

	user := ctx.Value("user").(*auth.User)

	token, exp, errGenerate := auth.GenerateToken(user)
	if errGenerate != nil {
		return nil, errGenerate
	}

	rsp := &TokenResponse{
		Token:      string(token),
		Expiration: exp.Truncate(time.Second),
	}

	return rsp, nil

}
