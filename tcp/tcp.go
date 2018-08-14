package tcp

import (
	"bytes"
	"fmt"
	"net"
	"posam/util/concurrentmap"
	"time"
)

type Connectivitier interface {
	Connect(string, string, time.Duration) (*net.TCPConn, error)
}

type Connectivity struct {
}

func (c *Connectivity) Connect(network string, address string, timeout time.Duration) (conn *net.TCPConn, err error) {
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
	Connection() *net.TCPConn
}

type TCPClient struct {
	Connectivitier
	ServerNetwork     string
	ServerAddress     string
	ServerTimeout     time.Duration
}

var connectionMap *concurrentmap.ConcurrentMap

func init() {
	connectionMap = concurrentmap.NewConcurrentMap()
}

func (t *TCPClient) Connection() *net.TCPConn {
	if connection, ok := connectionMap.Get(t.ServerAddress); ok {
		log.Println("INSTANCE FOUND")
		return connection.(*net.TCPConn)
	}
	if t.Connectivitier == nil {
		t.Connectivitier = &Connectivity{}
	}
	conn, err := t.Connect(t.ServerNetwork, t.ServerAddress, t.ServerTimeout)
	if err != nil {
		return conn
	}
	return t.addInstance(conn)
}

func (t *TCPClient) addConnection(conn *net.TCPConn) *net.TCPConn {
	log.Println("INSTANCE ADDED")
	connection := connectionMap.Set(t.ServerAddress, conn)
	return connection.(*net.TCPConn)
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
