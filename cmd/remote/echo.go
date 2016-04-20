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

	instance, username, password, errCredentials := getInstanceUsernamePassword(c)
	if errCredentials != nil {
		fmt.Printf("Error: %s\n", errCredentials.Error())
		os.Exit(1)
	}

	authTransport := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	authUser := grpc.WithPerRPCCredentials(&BasicAuthCredentials{Username: username, Password: password})
	conn, err := grpc.Dial(instance, authTransport, authUser)
	if err != nil {
		fmt.Printf("Error connecting to remote (%s): %s\n", instance, err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)
	msgValue := strings.Join(c.Args(), " ")

	if msg, err := client.Echo(context.Background(), &pb.EchoMessage{Value: msgValue}); err == nil {
		fmt.Println(msg.Value)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

}
