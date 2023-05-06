package yeelight_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/Gprisco/yeelightcontrol/yeelight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPacketConn struct {
	mock.Mock
	ReceiverAddr net.Addr
	Message      string
}

func (m *mockPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	args := m.Called(p)
	copy(p, []byte(args.String(0)))
	return 0, nil, args.Error(1)
}

func (m *mockPacketConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	m.Message = string(p)
	m.ReceiverAddr = addr

	args := m.Called(p, addr)
	return 0, args.Error(0)
}

func (*mockPacketConn) Close() error { return nil }

func (*mockPacketConn) LocalAddr() net.Addr { return nil }

func (*mockPacketConn) SetDeadline(t time.Time) error { return nil }

func (*mockPacketConn) SetReadDeadline(t time.Time) error { return nil }

func (*mockPacketConn) SetWriteDeadline(t time.Time) error { return nil }

const searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
	"MAN: \"ssdp:discover\"\r\n" +
	"ST: wifi_bulb\r\n\r\n"

func TestSendSearchMessage(t *testing.T) {
	// given a connection and an address on which to send the message
	wantedAddr := "1.2.3.4:1234"
	mockConn := &mockPacketConn{}

	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(nil)

	// when sending the message
	err := yeelight.SendSearchMessage(mockConn, wantedAddr)

	// then it should have sent the search message on the given address
	assert.Nil(t, err)
	assert.Equal(t, wantedAddr, mockConn.ReceiverAddr.String())
	assert.Equal(t, searchMessage, mockConn.Message)
}

func TestFailingAddressParsingInSendSearchMessage(t *testing.T) {
	// given a connection and an invalid address on which to send the message
	wantedAddr := "this is not a valid address"
	mockConn := &mockPacketConn{}

	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(nil)

	// when sending the message
	err := yeelight.SendSearchMessage(mockConn, wantedAddr)

	// then there's an error
	assert.Equal(t, &net.AddrError{Err: "missing port in address", Addr: "this is not a valid address"}, err)
}

func TestFailingWriteToSendSearchMessage(t *testing.T) {
	// given a connection and an address on which to send the message
	wantedAddr := "1.2.3.4:1234"
	mockConn := &mockPacketConn{}

	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(errors.New("mock failure when writing"))

	// when sending the message
	err := yeelight.SendSearchMessage(mockConn, wantedAddr)

	// then it should have sent the search message on the given address
	assert.Equal(t, errors.New("mock failure when writing"), err)
}

func TestReceiveResponse(t *testing.T) {
	// given a mock connection
	mockConn := &mockPacketConn{}
	mockConn.On("ReadFrom", mock.Anything).Return("mock response", nil)

	// when reading the response
	response, err := yeelight.ReceiveSearchResponse(mockConn)

	// then it should have received a response
	expectedResponse := make([]byte, 1024)
	copy(expectedResponse, "mock response")

	assert.Nil(t, err)
	assert.Equal(t, string(expectedResponse), response)
}

func TestFailingReceiveResponse(t *testing.T) {
	// given a mock connection
	mockConn := &mockPacketConn{}
	mockConn.On("ReadFrom", mock.Anything).Return("", errors.New("mock error"))

	// when reading the response
	_, err := yeelight.ReceiveSearchResponse(mockConn)

	// then it should have received an error
	assert.Equal(t, errors.New("mock error"), err)
}
