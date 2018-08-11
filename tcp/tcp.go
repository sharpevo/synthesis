package tcp

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type Connectioner interface {
	Connect(string, string, time.Duration) (net.Conn, error)
}

type Connection struct {
}

func (c *Connection) Connect(network string, address string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = net.Dial(network, address)
	if err != nil {
		return conn, err
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	conn.SetDeadline(time.Now().Add(timeout))
	return conn, nil
}

type TCPClienter interface {
	Send([]byte, []byte) ([]byte, error)
}

type TCPClient struct {
	Connectioner
	ServerNetwork string
	ServerAddress string
	ServerTimeout time.Duration
}

func (t *TCPClient) Send(message []byte, expected []byte) (resp []byte, err error) {
	conn, err := t.Connect(
		t.ServerNetwork,
		t.ServerAddress,
		t.ServerTimeout,
	)
	if err != nil {
		return resp, err
	}
	defer conn.Close()
	conn.Write(message)
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	resp = buf[:n]
	if !bytes.Equal(expected, resp) {
		return resp, fmt.Errorf("response error %x (%x)",
			resp,
			expected,
		)
	}
	return resp, nil
}
