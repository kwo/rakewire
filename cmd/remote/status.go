package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"os"
	"rakewire/api/pb"
	"time"
)

// Status retrieves the status of a remote instance
func Status(c *cli.Context) error {

	conn, errConnect := connect(c)
	if errConnect != nil {
		fmt.Printf("Error: %s\n", errConnect.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewStatusServiceClient(conn)

	if status, err := client.GetStatus(context.Background(), &pb.StatusRequest{}); err == nil {

		// TODO: better formatting, include days

		duration := time.Now().Truncate(time.Second).Sub(time.Unix(status.AppStart, 0))
		fmt.Printf("uptime: %s\n", duration.String())

		if len(status.Version) != 0 {
			fmt.Printf("version: %s\n", status.Version)
		}

		if status.BuildTime != 0 {
			fmt.Printf("build time: %s\n", time.Unix(status.BuildTime, 0).UTC().Format(time.RFC3339))
		}

		if len(status.BuildHash) != 0 {
			fmt.Printf("build hash: %s\n", status.BuildHash)
		}

	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
