package tcp

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type TCPClienter interface {
	Send([]byte, []byte) ([]byte, error)
}

type TCPClient struct {
	ServerNetwork string
	ServerAddress string
	ServerTimeout time.Duration
}

func (t *TCPClient) connect() (conn net.Conn, err error) {
	conn, err = net.Dial(t.ServerNetwork, t.ServerAddress)
	if t.ServerTimeout == 0 {
		t.ServerTimeout = 10 * time.Second
	}
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func (t *TCPClient) Send(message []byte, expected []byte) (resp []byte, err error) {
	conn, err := t.connect()
	conn.SetDeadline(time.Now().Add(t.ServerTimeout))
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
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
