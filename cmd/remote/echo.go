package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"os"
	"rakewire/api/pb"
	"strings"
)

// Echo tests the echo service
func Echo(c *cli.Context) {

	var endpoint string
	if c.NArg() > 0 {
		endpoint = c.Args()[0]
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	authTransport := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// TODO: retrieve username password
	authUser := grpc.WithPerRPCCredentials(&UsernamePasswordCredential{Username: "admin", Password: "123"})
	conn, err := grpc.Dial(endpoint, authTransport, authUser)
	if err != nil {
		fmt.Printf("Error connecting to remote (%s): %s\n", endpoint, err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)

	md := metadata.Pairs("auth", "123")
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)

	msgValue := strings.Join(c.Args()[1:], " ")
	if msg, err := client.Echo(ctx, &pb.EchoMessage{Value: msgValue}); err == nil {
		fmt.Println(msg.Value)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

}
