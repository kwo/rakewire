package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"rakewire/api/pb"
	"strings"
)

// Echo tests the echo service
func Echo(c *cli.Context) {

	var endpoint string
	var hostname string
	if c.NArg() > 0 {
		endpoint = c.Args()[0]
		hostname = strings.Split(endpoint, ":")[0]
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, hostname))}
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		fmt.Printf("Error connecting to remote (%s): %s\n", endpoint, err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)

	msgValue := strings.Join(c.Args()[1:], " ")
	if msg, err := client.Echo(context.Background(), &pb.EchoMessage{Value: msgValue}); err == nil {
		fmt.Println(msg.Value)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

}
