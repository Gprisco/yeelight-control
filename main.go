package main

import (
	"log"
	"net"
	"os"

	"github.com/Gprisco/yeelightcontrol/yeelight"
	"github.com/urfave/cli/v2"
)

const (
	// The multicast port on which to send the search message https://www.yeelight.com/en_US/developer
	multicastAddr = "239.255.255.250:1982"
	localAddr     = "0.0.0.0:0"
)

func handleSearchBulbs() {
	// Step 1: Create a UDP connection and bind to a local address.
	conn, err := net.ListenPacket("udp", localAddr)
	if err != nil {
		log.Printf("Failed to bind to local address: %v\n", err)
		return
	}
	defer conn.Close()

	bulbSearcher := yeelight.NewBulbSearcher(conn, multicastAddr)

	response, err := bulbSearcher.Search()

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
					handleSearchBulbs()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
