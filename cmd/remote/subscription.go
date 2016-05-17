package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/api/msg"
	"os"
	"strings"
	"time"
)

// SubscriptionAdd adds a new subscription
func SubscriptionAdd(c *cli.Context) error {

	var url string
	var groups []string
	var title string
	if c.NArg() >= 2 {
		url = c.Args()[0]
		groups = strings.Split(c.Args()[1], ",")
		if c.NArg() >= 3 {
			title = c.Args()[2]
		}
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	req := &msg.SubscriptionAddUpdateRequest{
		AddGroups: c.Bool("groups"),
		Subscription: &msg.Subscription{
			URL:      url,
			Groups:   groups,
			Title:    title,
			AutoRead: c.Bool("autoread"),
			AutoStar: c.Bool("autostar"),
		},
	}

	rsp := &msg.SubscriptionAddUpdateResponse{}

	if err := makeRequest(c, "subscriptions/add", req, rsp); err == nil {
		if len(rsp.Message) > 0 {
			fmt.Printf("%s: %s\n", msg.StatusText(rsp.Status), rsp.Message)
		} else {
			fmt.Println(msg.StatusText(rsp.Status))
		}
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}

// SubscriptionList retrieves the list of user subscriptions from the remote instance
func SubscriptionList(c *cli.Context) error {

	req := &msg.SubscriptionListRequest{}
	rsp := &msg.SubscriptionListResponse{}

	if c.NArg() == 1 {
		req.Filter = c.Args()[0]
	}

	if err := makeRequest(c, "subscriptions/list", req, rsp); err == nil {

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

		fmt.Printf("%-15s %s %s %-25s %-80s %-20s\n", "groups", "r", "s", "title", "url", "added")
		for _, sub := range rsp.Subscriptions {
			fmt.Printf("%-15s %s %s %-25s %-80s %-20s\n", strings.Join(sub.Groups, ", "), fmtBool(sub.AutoRead, "#"), fmtBool(sub.AutoStar, "*"), sub.Title, sub.URL, sub.Added.Format(time.RFC3339))
		}

	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}

// SubscriptionUnsubscribe unsubscribes from a subscription
func SubscriptionUnsubscribe(c *cli.Context) error {

	var url string
	if c.NArg() >= 1 {
		url = c.Args()[0]
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	req := &msg.UnsubscribeRequest{
		URL: url,
	}
	rsp := &msg.UnsubscribeResponse{}

	if err := makeRequest(c, "subscriptions/unsubscribe", req, rsp); err == nil {
		if len(rsp.Message) > 0 {
			fmt.Printf("%s: %s\n", msg.StatusText(rsp.Status), rsp.Message)
		} else {
			fmt.Println(msg.StatusText(rsp.Status))
		}
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
