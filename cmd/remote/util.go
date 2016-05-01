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
	errMissingHost     = errors.New("Missing host/port, try setting the --host option")
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
	auth := base64.StdEncoding.EncodeToString([]byte(z.Username + ":" + z.Password))
	return map[string]string{
		"authorization": "Basic " + auth,
	}, nil
}

// RequireTransportSecurity is part of the Credential interface,
func (z *BasicAuthCredentials) RequireTransportSecurity() bool {
	return true
}

// TokenCredentials implements Credentials for JWT authentication.
type TokenCredentials struct {
	Token string
}

// GetRequestMetadata is part of the Credential interface,
func (z *TokenCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + z.Token,
	}, nil
}

// RequireTransportSecurity is part of the Credential interface,
func (z *TokenCredentials) RequireTransportSecurity() bool {
	return true
}

func connect(c *cli.Context) (*grpc.ClientConn, error) {

	insecureSkipVerify := c.Parent().Bool("insecure")

	addr, username, password, token, errCredentials := getHostUsernamePasswordToken(c)
	if errCredentials != nil {
		return nil, errCredentials
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}

	var creds credentials.Credentials
	if len(token) == 0 {
		creds = &BasicAuthCredentials{Username: username, Password: password}
	} else {
		creds = &TokenCredentials{Token: token}
	}

	authTransport := grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	authUser := grpc.WithPerRPCCredentials(creds)
	return grpc.Dial(addr, authTransport, authUser)

}

func getHostUsernamePasswordToken(c *cli.Context) (host, username, password, token string, err error) {

	host = c.Parent().String("host")
	username = c.Parent().String("username")
	password = c.Parent().String("password")
	token = c.Parent().String("token")

	if len(host) == 0 {
		err = errMissingHost
		return
	}

	if len(token) == 0 {

		if len(username) == 0 {
			err = errMissingUsername
			return
		}

		if len(password) == 0 {
			err = errMissingPassword
			return
		}

	}

	return

}
