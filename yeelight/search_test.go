package yeelight

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----- Creating a PacketConn mock
type mockPacketConn struct {
	mock.Mock
	ReceiverAddr net.Addr
	Message      string
}

// ----- Implementing PacketConn interface
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

func TestSearch(t *testing.T) {
	// given a mock packet connection
	mockConn := &mockPacketConn{}
	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(nil)
	mockConn.On("ReadFrom", mock.Anything).Return("mock response", nil)

	// and a mock bulb searcher
	bulbSearcher := &BulbSearcher{
		conn:          mockConn,
		multicastAddr: "1.2.3.4:1234",
	}

	// when calling HandleSearch
	res, err := bulbSearcher.Search()

	// then err should be nil
	p := make([]byte, 1024)
	copy(p, "mock response")

	assert.Nil(t, err)
	assert.Equal(t, string(p), res)
}

func TestFailingWriteToSearch(t *testing.T) {
	// given a mock failing packet connection
	mockConn := &mockPacketConn{}
	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(errors.New("mock WriteTo error"))
	mockConn.On("ReadFrom", mock.Anything).Return("mock response", nil)

	// and a mock bulb searcher
	bulbSearcher := &BulbSearcher{
		conn:          mockConn,
		multicastAddr: "1.2.3.4:1234",
	}

	// when calling HandleSearch
	_, err := bulbSearcher.Search()

	// then err should be nil
	assert.Equal(t, errors.New("mock WriteTo error"), err)
}

func TestFailingReadFromSearch(t *testing.T) {
	// given a mock failing packet connection
	mockConn := &mockPacketConn{}
	mockConn.On("WriteTo", mock.Anything, mock.Anything).Return(nil)
	mockConn.On("ReadFrom", mock.Anything).Return("", errors.New("mock ReadFrom error"))

	// and a mock bulb searcher
	bulbSearcher := &BulbSearcher{
		conn:          mockConn,
		multicastAddr: "1.2.3.4:1234",
	}

	// when calling HandleSearch
	_, err := bulbSearcher.Search()

	// then err should be nil
	assert.Equal(t, errors.New("mock ReadFrom error"), err)
}
