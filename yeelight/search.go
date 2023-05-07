package yeelight

import (
	"net"
)

type Searcher interface {
	HandleSearch() (string, error)
}

type BulbSearcher struct {
	conn          net.PacketConn
	multicastAddr string
}

func NewBulbSearcher(conn net.PacketConn, multicastAddr string) *BulbSearcher {
	return &BulbSearcher{
		conn:          conn,
		multicastAddr: multicastAddr,
	}
}

const searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
	"MAN: \"ssdp:discover\"\r\n" +
	"ST: wifi_bulb\r\n\r\n"

func (s *BulbSearcher) sendSearchMessage(addr string) (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		return
	}

	_, err = s.conn.WriteTo([]byte(searchMessage), udpAddr)
	return
}

func (s *BulbSearcher) receiveSearchResponse() (string, error) {
	buf := make([]byte, 1024)
	for {
		_, _, err := s.conn.ReadFrom(buf)

		return string(buf), err
	}
}

func (s *BulbSearcher) Search() (res string, err error) {
	if err = s.sendSearchMessage(s.multicastAddr); err != nil {
		return
	}

	res, err = s.receiveSearchResponse()
	return
}
