package tcp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"synthesis/util/blockingqueue"
	"synthesis/util/concurrentmap"
	"time"
)

var clientMap *concurrentmap.ConcurrentMap

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	clientMap = concurrentmap.NewConcurrentMap()
}

func Instance(key string) *Client {
	if key == "" {
		for item := range clientMap.Iter() {
			return item.Value.(*Client)
		}
	} else {
		if clienti, ok := clientMap.Get(key); ok {
			return clienti.(*Client)
		} else {
			return nil
		}
	}
	return nil
}

func addInstance(client *Client) (*Client, bool) {
	if c, ok := clientMap.Get(client.Address); ok {
		return c.(*Client), true
	} else {
		clientMap.Set(client.Address, client)
		return client, false
	}
}

func ResetInstance() {
	for item := range clientMap.Iter() {
		client := item.Value.(*Client)
		log.Println("terminating client: ", client.Address)
		client.RequestQueue.Reset()
	}
	clientMap = concurrentmap.NewConcurrentMap()
}

type Clienter interface {
	connect() error
	Send([]byte, []byte) ([]byte, error)
}

type Client struct {
	Network      string
	Address      string
	Timeout      time.Duration
	RequestQueue *blockingqueue.BlockingQueue
	Conn         *net.TCPConn
}

func NewClient(
	network string,
	address string,
	seconds int,
) (*Client, error) {
	timeout := time.Duration(seconds) * time.Second
	client := &Client{
		Network:      network,
		Address:      address,
		Timeout:      timeout,
		RequestQueue: blockingqueue.NewBlockingQueue(),
	}
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	go client.launch()
	return client, nil
}

func (c *Client) connect() error {
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
	c.Conn = conn
	return nil
}

type Request struct {
	Message   []byte
	Expected  []byte
	Responsec chan Response
}

type Response struct {
	Message []byte
	Error   error
}

func (c *Client) launch() {
	//time.Sleep(3 * time.Second)
	log.Println("client launched")
	for {
		reqi, err := c.RequestQueue.Pop()
		if err != nil {
			log.Printf("client %q terminated\n", c.Address)
			return
		}
		req := reqi.(*Request)
		if c.Conn == nil {
			log.Println("connecting to the server...")
			err := c.connect()
			if err != nil {
				log.Println("failed to connect to the server:", err)
				req.Responsec <- Response{Error: err}
				continue
			}
		}
		c.send(req)
	}
}

func (c *Client) send(req *Request) {
	respc := req.Responsec
	resp := Response{}
	log.Println("sending request...", req.Message)
	c.Conn.SetDeadline(time.Now().Add(c.Timeout))
	n, err := c.Conn.Write(req.Message)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			log.Println("Reconnecting...")
			c.Conn = nil
		}
		resp.Error = err
		respc <- resp
		return
	}
	buf := make([]byte, 1536)
	c.Conn.SetDeadline(time.Now().Add(c.Timeout))
	n, err = c.Conn.Read(buf)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			log.Println("Reconnecting...")
			c.Conn = nil
		}
		resp.Error = err
		respc <- resp
		return
	}
	resp.Message = buf[:n]
	if !bytes.Equal(req.Expected, resp.Message) {
		resp.Error = fmt.Errorf("response error %v (%x)",
			resp.Message,
			req.Expected,
		)
	}
	respc <- resp
	log.Println("response received:", resp)
	return
}

func (c *Client) Send(
	message []byte,
	expected []byte,
) ([]byte, error) {
	req := Request{
		Message:   message,
		Expected:  expected,
		Responsec: make(chan Response),
	}
	c.RequestQueue.Push(&req)
	log.Println("waiting for response:", message)
	resp := <-req.Responsec
	return resp.Message, resp.Error
}
