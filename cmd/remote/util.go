package remote

import (
	"encoding/base64"
	"golang.org/x/net/context"
)

// UsernamePasswordCredential implements Credentials for username/password authentication.
type UsernamePasswordCredential struct {
	Username string
	Password string
}

// GetRequestMetadata is part of the Credential interface,
func (z *UsernamePasswordCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": z.makeAuthorizationToken(),
	}, nil
}

// RequireTransportSecurity is part of the Credential interface,
func (z *UsernamePasswordCredential) RequireTransportSecurity() bool {
	return true
}

// RequireTransportSecurity is part of the Credential interface,
func (z *UsernamePasswordCredential) makeAuthorizationToken() string {
	auth := base64.StdEncoding.EncodeToString([]byte(z.Username + ":" + z.Password))
	return "Basic " + auth
}
