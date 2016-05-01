package cmd

import (
	"github.com/codegangsta/cli"
)

// Version show version information and exit
func Version(c *cli.Context) error {
	showVersionInformation(c)
	return nil
}
