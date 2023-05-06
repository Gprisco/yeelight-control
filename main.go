package main

import (
	"log"
	"net"

	"github.com/Gprisco/yeelightcontrol/yeelight"
)

const (
	multicastAddr = "239.255.255.250:1982"
	localAddr     = "0.0.0.0:0"
)

func main() {
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
