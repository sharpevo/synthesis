package tcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"time"
)

type Connectivitier interface {
	Connect(string, string, time.Duration) (*net.TCPConn, error)
}

type Connectivity struct {
	Network      string
	Address      string
	Timeout      time.Duration
	RequestQueue *blockingqueue.BlockingQueue
	Conn         *net.TCPConn
}

var connectivityMap *concurrentmap.ConcurrentMap

func NewConnectivity(network string, address string, timeout time.Duration) *Connectivity {
	connectivity := &Connectivity{
		Network:      network,
		Address:      address,
		Timeout:      timeout,
		RequestQueue: blockingqueue.NewBlockingQueue(),
	}
	go connectivity.send()

	tcpAddr, _ := net.ResolveTCPAddr(network, address)
	return connectivityMap.Set(tcpAddr.String(), connectivity).(*Connectivity)
}

func (c *Connectivity) send() {
	time.Sleep(2 * time.Second)
	log.Printf("send launched")

	for {
		req := c.RequestQueue.Pop().(TCPRequest)

		respc := req.Respc
		resp := TCPResponse{}

		if c.Conn == nil {
			log.Println("connecting to the server...")
			err := c.connect()
			if err != nil {
				log.Println("failed to connect to the server:", err)
				resp.Error = err
				respc <- resp
				continue
			}
		}

		log.Println("sending request...", req.Message)

		c.Conn.Write(req.Message)
		buf := make([]byte, 1536)
		n, err := c.Conn.Read(buf)
		if err != nil {
			log.Println(err)
			resp.Error = err
			respc <- resp
			continue
		}
		resp.Response = buf[:n]
		if !bytes.Equal(req.Expected, resp.Response) {
			resp.Error = fmt.Errorf("response error %v (%x)",
				resp.Response,
				req.Expected,
			)
		}
		respc <- resp

		log.Println("response received:", resp)
	}
}

func (c *Connectivity) connect() error {
	tcpAddr, err := net.ResolveTCPAddr(c.Network, c.Address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP(c.Network, nil, tcpAddr)
	if err != nil {
		return err
	}

	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}
	conn.SetDeadline(time.Now().Add(c.Timeout))
	c.Conn = conn
	return nil
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
	conn.SetDeadline(time.Now().Add(timeout))
	c.Conn = conn
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
