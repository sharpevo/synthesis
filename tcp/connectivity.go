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
