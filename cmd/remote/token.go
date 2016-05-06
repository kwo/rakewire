package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"rakewire/api"
	"time"
)

// Token retrieves a generated authentication token from the server
func Token(c *cli.Context) error {

	shellExport := c.Bool("export")

	req := &api.TokenRequest{}
	rsp := &api.TokenResponse{}

	if err := makeRequest(c, "token", req, rsp); err == nil {
		if shellExport {
			fmt.Printf("export RAKEWIRE_TOKEN=\"%s\"\n", rsp.Token)
		} else {
			fmt.Printf("token: %s\n", rsp.Token)
			fmt.Printf("expires: %s\n", rsp.Expiration.Format(time.RFC3339))
		}
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
