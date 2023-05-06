package main

import (
	"fmt"
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
		fmt.Printf("Failed to bind to local address: %v\n", err)
		return
	}
	defer conn.Close()

	// Step 2: Multicast the search request message.
	yeelight.SendSearchMessage(conn, multicastAddr)

	// Step 3: Wait for the response message.
	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFrom(buf)

		if err != nil {
			fmt.Printf("Failed to receive message: %v\n", err)
			return
		}
		fmt.Printf("Received message from %s: %s\n", addr.String(), string(buf[:n]))
	}
}
