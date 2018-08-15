package tcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"posam/util/concurrentmap"
	"time"
)

type TCPClienter interface {
	Send([]byte, []byte) ([]byte, error)
	Connection() *net.TCPConn
}

type TCPClient struct {
	Connectivitier
	ServerNetwork     string
	ServerAddress     string
	ServerTimeout     time.Duration
	ServerConcurrency bool
}

func NewTCPClient(network string, address string, second int, concurrency bool) *TCPClient {
	timeout := time.Duration(second) * time.Second
	return &TCPClient{
		Connectivitier:    NewConnectivity(network, address, timeout),
		ServerNetwork:     network,
		ServerAddress:     address,
		ServerTimeout:     timeout,
		ServerConcurrency: concurrency,
	}
}

var connectionMap *concurrentmap.ConcurrentMap

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	connectivityMap = concurrentmap.NewConcurrentMap()
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
		log.Println(err)
		return conn
	}
	log.Printf("INSTANCE CREATED: %#v\n", conn)
	if !t.ServerConcurrency {
		return t.addConnection(conn)
	} else {
		return conn
	}
	return t.addConnection(conn)
}

func (t *TCPClient) addConnection(conn *net.TCPConn) *net.TCPConn {
	log.Println("INSTANCE ADDED")
	connection := connectionMap.Set(t.ServerAddress, conn)
	return connection.(*net.TCPConn)
}

type TCPRequest struct {
	Message  []byte
	Expected []byte
	Respc    chan TCPResponse
}

type TCPResponse struct {
	Response []byte
	Error    error
}

func (t *TCPClient) Send(message []byte, expected []byte) (resp []byte, err error) {
	req := TCPRequest{
		Message:  message,
		Expected: expected,
		Respc:    make(chan TCPResponse),
	}

	//tcpAddr, err := net.ResolveTCPAddr(t.ServerNetwork, t.ServerAddress)
	//if err != nil {
	//return resp, err
	//}

	//cnt, found := connectivityMap.Get(tcpAddr.String())
	//if !found {
	//fmt.Println("============ not found")
	//cnt = NewConnectivity(t.ServerNetwork, t.ServerAddress, t.ServerTimeout)
	//} else {
	//fmt.Println("============ found")
	//}
	//connectivity := cnt.(*Connectivity)
	connectivity := t.Connectivitier.(*Connectivity)

	//connectivity := connectivityMap.GetConnectivity(t.ServerNetwork, t.ServerAddress, t.ServerTimeout)
	//log.Println(connectivity)
	connectivity.RequestQueue.Push(req)

	//log.Println("pushed", req.Message)
	//t.Connectivitier.(Connectivity).RequestQueue.Push(req)
	log.Println("wait respc", message)
	tcpResponse := <-req.Respc
	log.Println("get respc", tcpResponse.Response)
	return tcpResponse.Response, tcpResponse.Error
}

func (t *TCPClient) send(message []byte, expected []byte) (resp []byte, err error) {
	conn := t.Connection()
	if t.ServerConcurrency {
		defer conn.Close()
		defer log.Println("connection closed")
	}

	conn.Write(message)
	buf := make([]byte, 1536)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = buf[:n]
	if !bytes.Equal(expected, resp) {
		return resp, fmt.Errorf("response error %v (%x)",
			resp,
			expected,
		)
	}
	return resp, nil
}
