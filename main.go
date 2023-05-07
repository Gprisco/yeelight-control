package main

import (
	"log"
	"net"
	"os"

	"github.com/Gprisco/yeelightcontrol/yeelight"
	"github.com/urfave/cli/v2"
)

const (
	multicastAddr = "239.255.255.250:1982"
	localAddr     = "0.0.0.0:0"
)

func searchBulbs() {
	// Step 1: Create a UDP connection and bind to a local address.
	conn, err := net.ListenPacket("udp", localAddr)
	if err != nil {
		log.Printf("Failed to bind to local address: %v\n", err)
		return
	}
	defer conn.Close()

	// Step 2: Multicast the search request message.
	if err = yeelight.SendSearchMessage(conn, multicastAddr); err != nil {
		log.Printf("Failed to broadcast search message: %v\n", err)
		return
	}

	// Step 3: Wait for the response message.
	response, err := yeelight.ReceiveSearchResponse(conn)

	if err != nil {
		log.Printf("Failed to get response from yeelight: %v\n", err)
		return
	}

	log.Println(response)
}

func main() {
	app := &cli.App{
		Name:  "yeelightctl",
		Usage: "Control yeelight bulbs under your local network",
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "Search bulbs on local network",
				Action: func(cCtx *cli.Context) error {
					searchBulbs()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
