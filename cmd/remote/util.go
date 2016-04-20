package remote

import (
	"encoding/base64"
	"errors"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
)

var (
	errMissingInstance = errors.New("Missing instance, try setting the --instance option")
	errMissingUsername = errors.New("Missing username, try setting the --username option")
	errMissingPassword = errors.New("Missing password, try setting the --password option")
)

func getInstanceUsernamePassword(c *cli.Context) (instance, username, password string, err error) {

	instance = c.Parent().String("instance")
	username = c.Parent().String("username")
	password = c.Parent().String("password")

	if instance == "" {
		err = errMissingInstance
		return
	}

	if username == "" {
		err = errMissingUsername
		return
	}

	if password == "" {
		err = errMissingPassword
		return
	}

	return

}

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
