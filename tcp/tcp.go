package tcp

import (
	"bytes"
	"fmt"
	"net"
	"posam/util/concurrentmap"
	"time"
)

type Connectioner interface {
	Connect(string, string, time.Duration) (*net.TCPConn, error)
}

type Connection struct {
}

func (c *Connection) Connect(network string, address string, timeout time.Duration) (conn *net.TCPConn, err error) {
	//conn, err := net.Dial(network, address)

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return conn, err
	}

	conn, err = net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		return conn, err
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	//conn.SetDeadline(time.Now().Add(timeout))
	return conn, nil
}

type TCPClienter interface {
	Send([]byte, []byte) ([]byte, error)
	Instance() *net.TCPConn
}

type TCPClient struct {
	Connectioner
	ServerNetwork string
	ServerAddress string
	ServerTimeout time.Duration
}

var instanceMap *concurrentmap.ConcurrentMap

func init() {
	instanceMap = concurrentmap.NewConcurrentMap()
}

func (t *TCPClient) Instance() *net.TCPConn {
	fmt.Println("INSTANCE")
	if instance, ok := instanceMap.Get(t.ServerAddress); ok {
		fmt.Println("INSTANCE FOUND")
		return instance.(*net.TCPConn)
	}
	if t.Connectioner == nil {
		t.Connectioner = &Connection{}
	}
	conn, err := t.Connect(t.ServerNetwork, t.ServerAddress, t.ServerTimeout)
	if err != nil {
		return conn
	}
	return t.addInstance(conn)
}

func (t *TCPClient) addInstance(conn *net.TCPConn) *net.TCPConn {
	fmt.Println("INSTANCE ADDED")
	instance := instanceMap.Set(t.ServerAddress, conn)
	return instance.(*net.TCPConn)
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
		return resp, fmt.Errorf("response error %v (%x)",
			resp,
			expected,
		)
	}
	return resp, nil
}
