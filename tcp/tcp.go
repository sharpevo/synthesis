package tcp

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type TCPClient struct {
	ServerNetwork string
	ServerAddress string
}

func (t *TCPClient) connect() (conn net.Conn, err error) {
	conn, err = net.Dial(t.ServerNetwork, t.ServerAddress)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func (t *TCPClient) Send(message []byte, expected []byte) (resp []byte, err error) {
	conn, err := t.connect()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
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
