package tcp

import (
	"bufio"
	"fmt"
	"net"
)

type TCP struct {
	Network string
	Address string
}

func (t *TCP) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial(t.Network, t.Address)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func (t *TCP) SendString(message string) string {
	conn, err := t.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(conn, message)
	response, err := bufio.NewReader(conn).ReadString('\n')
	return response
}

func (t *TCP) SendByte(message []byte) []byte {
	conn, err := t.Connect()
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(message)
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	return buf[:n]
}
