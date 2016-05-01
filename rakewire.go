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
				cli.BoolFlag{
					Name:   "k, insecure",
					EnvVar: "RAKEWIRE_INSECURE",
					Usage:  "Skip verification of TLS certificate, use with a self-signed certificate",
				},
				cli.StringFlag{
					Name:   "bind",
					Value:  "0.0.0.0:8888",
					EnvVar: "RAKEWIRE_BIND",
					Usage:  "host:port on which httpd should listen, defaults to 0.0.0.0:8888",
				},
				cli.StringFlag{
					Name:   "host",
					Value:  "localhost:8888",
					EnvVar: "RAKEWIRE_HOST",
					Usage:  "host:port on which httpd will be (publicly) accessible, defaults to localhost:8888",
				},
				cli.StringFlag{
					Name:   "tlscert",
					Value:  "rakewire.crt",
					EnvVar: "RAKEWIRE_TLSCERT",
					Usage:  "TLS certificate file",
				},
				cli.StringFlag{
					Name:   "tlskey",
					Value:  "rakewire.key",
					EnvVar: "RAKEWIRE_TLSKEY",
					Usage:  "TLS key file",
				},
				cli.StringFlag{
					Name:   "useragent",
					Value:  strings.TrimSpace("Rakewire " + Version),
					EnvVar: "RAKEWIRE_USERAGENT",
					Usage:  "useragent string",
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
			Usage: "perform a data integrity check on the database",
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
					Name:   "host",
					Value:  "localhost",
					EnvVar: "RAKEWIRE_HOST",
					Usage:  "Comma-separated hostnames and IPs to generate a certificate for, :port will be removed",
				},
				cli.IntFlag{
					Name:  "bits",
					Value: 2048,
					Usage: "Size of RSA key to generate. Ignored if --curve is set",
				},
				cli.StringFlag{
					Name:  "curve",
					Usage: "ECDSA curve to use to generate a key. Valid values are P224, P256, P384 and P521",
				},
				cli.StringFlag{
					Name:   "tlscert",
					Value:  "rakewire.crt",
					EnvVar: "RAKEWIRE_TLSCERT",
					Usage:  "Location of the certificate file to generate. Will not overwrite existing file.",
				},
				cli.StringFlag{
					Name:   "tlskey",
					Value:  "rakewire.key",
					EnvVar: "RAKEWIRE_TLSKEY",
					Usage:  "Location of the key file to generate. Will not overwrite existing file.",
				},
			},
		},
		{
			Name:      "useradd",
			Usage:     "add user",
			ArgsUsage: "<username> <password> [role[,role]]",
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
				cli.BoolFlag{
					Name:   "k, insecure",
					EnvVar: "RAKEWIRE_INSECURE",
					Usage:  "Skip verification of remote instance TLS certificate",
				},
				cli.StringFlag{
					Name:   "host",
					Value:  "localhost:8888",
					EnvVar: "RAKEWIRE_HOST",
					Usage:  "fqdn:port of the remote server, defaults to localhost:8888",
				},
				cli.StringFlag{
					Name:   "username",
					EnvVar: "RAKEWIRE_USERNAME",
					Usage:  "name of rakewire user",
				},
				cli.StringFlag{
					Name:   "password",
					EnvVar: "RAKEWIRE_PASSWORD",
					Usage:  "password for the rakewire user",
				},
				cli.StringFlag{
					Name:   "token",
					EnvVar: "RAKEWIRE_TOKEN",
					Usage:  "jwt authentication as alternative to username/password",
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
				{
					Name:   "token",
					Usage:  "get authentication token",
					Action: remote.Token,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "x, export",
							Usage: "generate shell command to store token in environment variable",
						},
					},
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
