package remote

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/api/msg"
)

// Token retrieves a generated authentication token from the server
func Token(c *cli.Context) error {

	shellExport := c.Bool("export")

	req := &msg.TokenRequest{}
	rsp := &msg.TokenResponse{}

	if err := makeRequest(c, "token", req, rsp); err == nil {
		if shellExport {
			fmt.Printf("export RAKEWIRE_TOKEN=\"%s\"\n", rsp.Token)
		} else {
			fmt.Printf("token: %s\n", rsp.Token)
			fmt.Printf("expires: %s\n", time.Unix(rsp.Expiration, 0).Format(time.RFC3339))
		}
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
