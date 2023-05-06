package yeelight

import (
	"net"
)

const searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
	"MAN: \"ssdp:discover\"\r\n" +
	"ST: wifi_bulb\r\n\r\n"

func SendSearchMessage(conn net.PacketConn, addr string) (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return
	}

	_, err = conn.WriteTo([]byte(searchMessage), udpAddr)
	return
}

func ReceiveSearchResponse(conn net.PacketConn) (string, error) {
	buf := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFrom(buf)

		return string(buf), err
	}
}
