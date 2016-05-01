package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"os"
	"rakewire/api/pb"
	"time"
)

// Token retrieves a generated authentication token from the server
func Token(c *cli.Context) error {

	shellExport := c.Bool("export")

	conn, errConnect := connect(c)
	if errConnect != nil {
		fmt.Printf("Error: %s\n", errConnect.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewTokenServiceClient(conn)

	if rsp, err := client.GetToken(context.Background(), &pb.TokenRequest{}); err == nil {
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
