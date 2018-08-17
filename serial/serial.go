package serial

import (
	"bytes"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"posam/protocol/modbus"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
)

var clientMap *concurrentmap.ConcurrentMap

func init() {
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
	key := string(client.Name)
	if c, ok := clientMap.Get(key); ok {
		return c.(*Client), true
	} else {
		clientMap.Set(key, client)
		return client, false
	}
}

func ResetInstance() {
	for item := range clientMap.Iter() {
		client := item.Value.(*Client)
		log.Println("terminating client: ", client.Name)
		client.RequestQueue.Reset()
	}
	clientMap = concurrentmap.NewConcurrentMap()
}

type Clienter interface {
	connect() error
	Send([]byte, []byte) ([]byte, error)
}

type Client struct {
	Name     string
	BaudRate int

	DataBits int
	StopBits int
	Parity   int

	Connection   *serial.Port
	RequestQueue *blockingqueue.BlockingQueue
}

func NewClient(
	name string,
	baud int,
	databits int,
	stopbits int,
	parity int,
) (*Client, error) {
	client := &Client{
		Name:         name,
		BaudRate:     baud,
		DataBits:     databits,
		StopBits:     stopbits,
		Parity:       parity,
		RequestQueue: blockingqueue.NewBlockingQueue(),
	}
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	go client.launch()
	return client, nil
}

func (c *Client) connect() error {
	log.Printf("Opening serial port %q...", c.Name)
	conf := &serial.Config{
		Name: c.Name,
		Baud: c.BaudRate,
	}
	openedPort, err := serial.OpenPort(conf)
	if err != nil {
		log.Println(err)
		return err
	}
	c.Connection = openedPort
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
	log.Println("serial client launched")
	for {
		reqi, err := c.RequestQueue.Pop()
		if err != nil {
			log.Printf("serial client %x terminated\n", c.Name)
			return
		}
		req := reqi.(*Request)
		if c.Connection == nil {
			log.Println("connecting to the the serial device...")
			err := c.connect()
			if err != nil {
				log.Println("failed to connect to the serial device:", err)
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

	_, err := c.Connection.Write(req.Message)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			log.Println("Reconnecting...")
			c.Connection = nil
		}
		resp.Error = err
		respc <- resp
		return
	}

	max := len(req.Expected)
	buf := make([]byte, max)
	cnt := 0
	for {
		n, err := c.Connection.Read(buf)
		if err != nil {
			resp.Error = err
			respc <- resp
			return
		}
		cnt += n
		resp.Message = append(resp.Message, buf[:n]...)
		if cnt >= max || n == 0 {
			break
		}
	}
	log.Printf("%x | %x", req.Expected, resp.Message)
	if !bytes.Equal(req.Expected, resp.Message) {
		resp.Error = fmt.Errorf(
			"invalid response code %x (%x)",
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
	modbus.AppendCRC(&message)
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

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
