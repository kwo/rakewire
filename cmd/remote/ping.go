package remote

import (
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"io"
	"os"
	"os/signal"
	"rakewire/api/pb"
	"syscall"
	"time"
)

// Ping retrieves the pings of a remote instance
func Ping(c *cli.Context) error {

	conn, errConnect := connect(c)
	if errConnect != nil {
		fmt.Printf("Error: %s\n", errConnect.Error())
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)

	if stream, errRequest := client.Ping(context.Background(), &pb.PingRequest{}); errRequest == nil {

		killsignal := make(chan os.Signal, 1)
		signal.Notify(killsignal, syscall.SIGINT, syscall.SIGTERM)
		ticker := time.NewTicker(time.Millisecond * 100) // must be smaller than interval from server
		done := stream.Context().Done()

	listen:
		for {
			if rsp, err := stream.Recv(); err == nil {
				fmt.Printf("ping: %s\n", time.Unix(rsp.Time, 0).Format(time.RFC3339))
			} else if err == io.EOF {
				fmt.Println(err.Error())
				break listen
			} else {
				fmt.Printf("Error: %s\n", err.Error())
				break listen
			}
			select {
			case <-ticker.C:
			case <-done:
				break listen
			case <-killsignal:
				break listen
			}
		} // loop
		fmt.Println("exiting...")
		ticker.Stop()

	} else {
		fmt.Printf("Error: %s\n", errRequest.Error())
		os.Exit(1)
	}

	return nil

}
