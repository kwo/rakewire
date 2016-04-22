package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"rakewire/cmd"
	"rakewire/cmd/remote"
	"strings"
)

// application level variables
var (
	Version   = ""
	BuildTime = ""
	BuildHash = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "Rakewire"
	app.Usage = "Feed Reader"
	app.HideVersion = true
	app.Version = strings.TrimSpace(fmt.Sprintf("%s %s %s", Version, BuildTime, BuildHash))
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "v, verbose",
			EnvVar: "RAKEWIRE_VERBOSE,VERBOSE",
			Usage:  "log more information to console",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "f, file",
					Value:  "rakewire.db",
					EnvVar: "RAKEWIRE_FILE",
					Usage:  "location of the database file",
				},
			},
			Action: cmd.Init,
		},
		{
			Name:  "start",
			Usage: "start rakewire",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "f, file",
					Value:  "rakewire.db",
					EnvVar: "RAKEWIRE_FILE",
					Usage:  "location of the database file",
				},
				cli.StringFlag{
					Name:   "p, pid",
					Value:  "rakewire.pid",
					EnvVar: "RAKEWIRE_PID",
					Usage:  "location of the pid file",
				},
				cli.IntFlag{
					Name:   "fetch.timeoutsecs",
					Value:  20,
					EnvVar: "RAKEWIRE_FETCH_TIMEOUTSECS",
					Usage:  "fetcher timeout",
				},
				cli.IntFlag{
					Name:   "fetch.workers",
					Value:  10,
					EnvVar: "RAKEWIRE_FETCH_WORKERS",
					Usage:  "fetcher workers",
				},
				cli.StringFlag{
					Name:   "fetch.useragent",
					Value:  strings.TrimSpace("Rakewire " + Version),
					EnvVar: "RAKEWIRE_FETCH_USERAGENT",
					Usage:  "fetcher useragent string",
				},
				cli.StringFlag{
					Name:   "httpd.address",
					Value:  "",
					EnvVar: "RAKEWIRE_HTTPD_ADDRESS",
					Usage:  "ip address on which httpd should listen",
				},
				cli.StringFlag{
					Name:   "httpd.host",
					Value:  "localhost",
					EnvVar: "RAKEWIRE_HTTPD_HOST",
					Usage:  "domain name at which httpd service can be reached",
				},
				cli.IntFlag{
					Name:   "httpd.port",
					Value:  8888,
					EnvVar: "RAKEWIRE_HTTPD_PORT",
					Usage:  "httpd port",
				},
				cli.StringFlag{
					Name:   "httpd.tlscertfile",
					EnvVar: "RAKEWIRE_HTTPD_TLSCERTFILE",
					Usage:  "TLS certificate file",
				},
				cli.StringFlag{
					Name:   "httpd.tlskeyfile",
					EnvVar: "RAKEWIRE_HTTPD_TLSKEYFILE",
					Usage:  "TLS key file",
				},
				cli.IntFlag{
					Name:   "poll.batchmax",
					Value:  10,
					EnvVar: "RAKEWIRE_POLL_BATCHMAX",
					Usage:  "maximum number of feeds to poll at once",
				},
				cli.IntFlag{
					Name:   "poll.intervalsecs",
					Value:  5,
					EnvVar: "RAKEWIRE_POLL_INTERVALSECS",
					Usage:  "how often to poll feeds",
				},
			},
			Action: cmd.Start,
		},
		{
			Name:  "check",
			Usage: "check database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "f, file",
					Value:  "rakewire.db",
					EnvVar: "RAKEWIRE_FILE",
					Usage:  "location of the database file",
				},
			},
			Action: cmd.Check,
		},
		{
			Name:   "certgen",
			Usage:  "generate tls certificate",
			Action: cmd.CertGen,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Value: "localhost",
					Usage: "Comma-separated hostnames and IPs to generate a certificate for",
				},
				cli.StringFlag{
					Name:  "start-date",
					Value: "",
					Usage: "Creation date formatted as Jan 1 15:04:05 2011",
				},
				cli.IntFlag{
					Name:  "duration-days",
					Value: 365,
					Usage: "Number of days that certificate is valid for",
				},
				cli.BoolFlag{
					Name:  "ca",
					Usage: "whether this cert should be its own Certificate Authority",
				},
				cli.IntFlag{
					Name:  "rsa-bits",
					Value: 2048,
					Usage: "Size of RSA key to generate. Ignored if --ecdsa-curve is set",
				},
				cli.StringFlag{
					Name:  "ecdsa-curve",
					Value: "",
					Usage: "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521",
				},
			},
		},
		{
			Name:      "useradd",
			Usage:     "add user",
			ArgsUsage: "<username> <roles>",
			Action:    cmd.UserAdd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "f, file",
					Value:  "rakewire.db",
					EnvVar: "RAKEWIRE_FILE",
					Usage:  "location of the database file",
				},
			},
		},
		{
			Name:    "remote",
			Aliases: []string{"r"},
			Usage:   "manage remote rakewire instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "i, instance",
					EnvVar: "RAKEWIRE_INSTANCE",
					Usage:  "name of remote rakewire instance in the format host:port",
				},
				cli.StringFlag{
					Name:   "u, username",
					EnvVar: "RAKEWIRE_USERNAME",
					Usage:  "name of rakewire user",
				},
				cli.StringFlag{
					Name:   "p, password",
					EnvVar: "RAKEWIRE_PASSWORD",
					Usage:  "password for the rakewire user",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "ping",
					Usage:  "get ping stream from server",
					Action: remote.Ping,
				},
				{
					Name:   "status",
					Usage:  "get instance status",
					Action: remote.Status,
				},
			},
		},
		{
			Name:   "version",
			Usage:  "print version and exit",
			Action: cmd.Version,
		},
	}
	app.Run(os.Args)
}
