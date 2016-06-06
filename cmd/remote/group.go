package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/api/msg"
	"os"
)

// GroupList retrieves the list of user groups from the remote instance
func GroupList(c *cli.Context) error {

	req := &msg.GroupListRequest{}
	rsp := &msg.GroupListResponse{}

	if err := makeRequest(c, "groups/list", req, rsp); err == nil {

		if rsp.Status != 0 {
			if len(rsp.Message) > 0 {
				fmt.Printf("%s: %s\n", msg.StatusText(rsp.Status), rsp.Message)
			} else {
				fmt.Println(msg.StatusText(rsp.Status))
			}
			return nil
		}

		fmt.Println("--- group listing ---")
		for _, group := range rsp.Groups {
			fmt.Printf("%-15s\n", group.Name)
		}

	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
