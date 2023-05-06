package yeelight

import (
	"fmt"
	"net"
)

const searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
	"MAN: \"ssdp:discover\"\r\n" +
	"ST: wifi_bulb\r\n\r\n"

func SendSearchMessage(conn net.PacketConn, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Failed to resolve multicast address: %v\n", err)
		return
	}
	if _, err := conn.WriteTo([]byte(searchMessage), udpAddr); err != nil {
		fmt.Printf("Failed to multicast message: %v\n", err)
		return
	}
}

func ReceiveSearchResponse(conn net.PacketConn) (string, error) {
	buf := make([]byte, 2048)
	for {
		_, _, err := conn.ReadFrom(buf)

		return string(buf), err
	}
}
