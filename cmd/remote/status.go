package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"github.com/kwo/rakewire/api"
	"time"
)

// Status retrieves the status of a remote instance
func Status(c *cli.Context) error {

	req := &api.StatusRequest{}
	status := &api.StatusResponse{}

	if err := makeRequest(c, "status", req, status); err == nil {

		// TODO: better formatting, include days

		duration := time.Now().Truncate(time.Second).Sub(status.AppStart)
		fmt.Printf("uptime: %s\n", duration.String())

		if len(status.Version) != 0 {
			fmt.Printf("version: %s\n", status.Version)
		}

		if !status.BuildTime.IsZero() {
			fmt.Printf("build time: %s\n", status.BuildTime.UTC().Format(time.RFC3339))
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
