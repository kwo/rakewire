package remote

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	errMissingInstance = errors.New("Missing host/port, try setting the --host/--port options")
	errMissingUsername = errors.New("Missing username, try setting the --username option")
	errMissingPassword = errors.New("Missing password, try setting the --password option")
)

// BasicAuthCredentials implements Credentials for username/password authentication.
type BasicAuthCredentials struct {
	Username string
	Password string
}

// GetRequestMetadata is part of the Credential interface,
func (z *BasicAuthCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": z.makeAuthorizationToken(),
	}, nil
}

// RequireTransportSecurity is part of the Credential interface,
func (z *BasicAuthCredentials) RequireTransportSecurity() bool {
	return true
}

// RequireTransportSecurity is part of the Credential interface,
func (z *BasicAuthCredentials) makeAuthorizationToken() string {
	auth := base64.StdEncoding.EncodeToString([]byte(z.Username + ":" + z.Password))
	return "Basic " + auth
}

func connect(c *cli.Context) (*grpc.ClientConn, error) {

	insecureSkipVerify := c.Parent().Bool("insecure")

	instance, username, password, errCredentials := getInstanceUsernamePassword(c)
	if errCredentials != nil {
		return nil, errCredentials
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}

	authTransport := grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	authUser := grpc.WithPerRPCCredentials(&BasicAuthCredentials{Username: username, Password: password})
	return grpc.Dial(instance, authTransport, authUser)

}

func getInstanceUsernamePassword(c *cli.Context) (instance, username, password string, err error) {

	instance = c.Parent().String("host")
	username = c.Parent().String("username")
	password = c.Parent().String("password")

	if len(instance) == 0 {
		err = errMissingInstance
		return
	}

	if len(username) == 0 {
		err = errMissingUsername
		return
	}

	if len(password) == 0 {
		err = errMissingPassword
		return
	}

	return

}
