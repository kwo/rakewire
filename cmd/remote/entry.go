package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/api/msg"
	"os"
	"time"
)

// EntryList retrieves the list of entries for a subscription
func EntryList(c *cli.Context) error {

	req := &msg.EntryListRequest{}
	rsp := &msg.EntryListResponse{}

	if c.NArg() == 1 {
		req.Subscription = c.Args()[0]
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	if err := makeRequest(c, "entries/list", req, rsp); err == nil {

		if rsp.Status != 0 {
			if len(rsp.Message) > 0 {
				fmt.Printf("%s: %s\n", msg.StatusText(rsp.Status), rsp.Message)
			} else {
				fmt.Println(msg.StatusText(rsp.Status))
			}
			return nil
		}

		fmtBool := func(value bool, marker string) string {
			if value {
				return marker
			}
			return " "
		}

		fmt.Printf("%s %s %-25s %-80s %-20s\n", "u", "s", "updated", "title", "guid")
		for _, entry := range rsp.Entries {
			fmt.Printf("%s %s %-25s %-80s %-20s\n", fmtBool(entry.Read, "#"), fmtBool(entry.Star, "*"), entry.Updated.Format(time.RFC3339), entry.Title, entry.GUID)
		}

	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
