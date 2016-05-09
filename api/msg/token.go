package msg

// TokenRequest defines the token request
type TokenRequest struct{}

// TokenResponse defines the token response
type TokenResponse struct {
	Token      string `json:"token,omitempty"`
	Expiration int64  `json:"expiration,omitempty"`
}
