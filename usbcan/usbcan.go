package usbcan

import (
	"controlcan"
	"encoding/hex"
	"fmt"
	"log"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"time"
)

const (
	STATUS_CODE_RECEIVED         = 0x00
	STATUS_CODE_COMPLETED        = 0x01
	STATUS_CODE_ERROR            = 0x02
	STATUS_CODE_MOTOR_INVALID    = 0x10
	STATUS_CODE_MAILBOX_FULL     = 0xE0
	STATUS_CODE_MAILBOX_OVERFLOW = 0xE1
)

var deviceMap *concurrentmap.ConcurrentMap
var clientMap *concurrentmap.ConcurrentMap

func init() {
	clientMap = concurrentmap.NewConcurrentMap()
	deviceMap = concurrentmap.NewConcurrentMap()
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
	key := string(client.DevID)
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
		log.Println("terminating client: ", client.DevID)
		client.RequestQueue.Reset()
	}
	clientMap = concurrentmap.NewConcurrentMap()
	deviceMap = concurrentmap.NewConcurrentMap()
}

type Clienter interface {
	connect() error
	Send([]byte, []byte, int, []byte, int) ([]byte, error)
}

type Client struct {
	DevType  int
	DevIndex int
	DevID    int // used as frame id

	CanIndex int

	AccCode int
	AccMask int
	Filter  int
	Timing0 int
	Timing1 int
	Mode    int

	RequestQueue       *blockingqueue.BlockingQueue
	ReceptionMap       *concurrentmap.ConcurrentMap
	InstructionCodeMap *concurrentmap.ConcurrentMap

	Sendable bool
}

func NewClient(
	devType int,
	devIndex int,
	devID int,
	canIndex int,
	accCode int,
	accMask int,
	filter int,
	timing0 int,
	timing1 int,
	mode int,
) (*Client, error) {
	client := &Client{
		DevType:  devType,
		DevIndex: devIndex,
		DevID:    devID,
		CanIndex: canIndex,
		AccCode:  accCode,
		AccMask:  accMask,
		Filter:   filter,
		Timing0:  timing0,
		Timing1:  timing1,
		Mode:     mode,
		Sendable: true,
	}
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	client.RequestQueue = blockingqueue.NewBlockingQueue()
	client.ReceptionMap = concurrentmap.NewConcurrentMap()
	client.InstructionCodeMap = concurrentmap.NewConcurrentMap()
	for index := range [256]byte{} {
		client.InstructionCodeMap.Set(
			hex.EncodeToString([]byte{byte(index)}),
			false,
		)
	}

	//err := client.connect()
	//if err != nil {
	//// TODO: notification
	//log.Println("Connect Error: ", err)
	//return client, err
	//}

	go client.launch()
	return client, nil
}

func (c *Client) deviceKey() string {
	return fmt.Sprintf("%v-%v", c.DevType, c.DevIndex)
}

func (c *Client) connect() (err error) {
	if _, ok := deviceMap.Get(c.deviceKey()); !ok {
		log.Printf(
			"Opening device(type %T, index %#v)...\n",
			c.DevType,
			c.DevIndex,
		)
		if err = controlcan.OpenDevice(
			int(c.DevType),
			int(c.DevIndex),
			0,
		); err != nil {
			return err
		}
		deviceMap.Set(c.deviceKey(), true)
	}
	config := controlcan.InitConfig{
		AccCode: c.AccCode,
		AccMask: c.AccMask,
		Filter:  c.Filter,
		Timing0: c.Timing0,
		Timing1: c.Timing1,
		Mode:    c.Mode,
	}
	log.Printf(
		"Initiating CAN(type %v, index %v, can %v)...\n",
		c.DevType,
		c.DevIndex,
		c.CanIndex,
	)
	if err = controlcan.InitCAN(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
		config,
	); err != nil {
		return err
	}
	log.Printf("Starting CAN(type %v, index %v)...\n", c.DevType, c.DevIndex)
	if err = controlcan.StartCAN(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
	); err != nil {
		return err
	}
	return nil
}

type Request struct {
	InstructionCode byte
	Message         []byte
	RecExpected     []byte
	ComExpected     []byte
	Responsec       chan Response
}

type Response struct {
	Message []byte
	Error   error
}

func (c *Client) launch() {
	log.Println("canalyst client launched")
	err := c.connect()
	if err != nil {
		// TODO: notification
		log.Println("Connect Error: ", err)
		return
	}
	go c.receive()
	for {
		reqi, err := c.RequestQueue.Pop()
		log.Println("processing", reqi)
		if err != nil {
			log.Printf("canalyst client sender %v terminated\n", c.deviceKey())
			return
		}
		req := reqi.(*Request)
		c.ReceptionMap.Set(
			hex.EncodeToString([]byte{req.InstructionCode}),
			req,
		)
		c.send(req)
	}
}

func (c *Client) receive() {
	log.Printf("listening %v...\n", c.deviceKey())
	for {
		time.Sleep(100 * time.Millisecond)
		pReceive := make([]controlcan.CanObj, controlcan.FRAME_LENGTH_OF_RECEPTION)
		count, err := controlcan.Receive(
			c.DevType,
			c.DevIndex,
			c.CanIndex,
			pReceive,
			100,
		)
		if err != nil || count < 0 {
			log.Printf("canalyst client receiver %v terminated\n", c.deviceKey())
			return
		}
		if count == 0 {
			continue
		}
		log.Printf("data received: %#v\n", pReceive[:count])
		for _, canObj := range pReceive[:count] {
			resp := Response{}
			// TODO: ? filter based on frame id
			req, err := c.findRequestByResponse(canObj.Data[:])
			if err != nil {
				log.Println(err)
				// TODO: notification
				continue
			}
			resp.Message = canObj.Data[:]
			req.Responsec <- resp
		}
	}
}

func (c *Client) parseFunctionCode(data []byte) (byte, error) {
	code := data[0]
	switch code {
	case 0xE0:
		c.Sendable = false
		return code, fmt.Errorf("mailbox is full (E0): %#v\n", data)
	case 0xE1:
		return code, fmt.Errorf("mailbox is overflow (E1): %#v\n", data)
	default:
	}
	return code, nil
}

func (c *Client) parseInstructionCode(data []byte) (code byte, err error) {
	code = data[7]
	if _, ok := c.InstructionCodeMap.Get(
		hex.EncodeToString([]byte{code}),
	); !ok {
		return code, fmt.Errorf("invalid instruction code")
	}
	return code, nil
}

func (c *Client) findRequestByResponse(data []byte) (request *Request, err error) {
	_, err = c.parseFunctionCode(data[:])
	if err != nil {
		return request, err
	}
	instCode, err := c.parseInstructionCode(data)
	if err != nil {
		return request, err
	}
	//for item := range c.RequestQueue.Iter() {
	log.Printf("parsing request: %s\n", c.ReceptionMap)
	for item := range c.ReceptionMap.Iter() {
		reqi := item.Value
		req, ok := reqi.(*Request)
		if !ok {
			err = fmt.Errorf("invalid request: %#v", reqi)
			log.Println(err)
			continue
		}
		log.Printf("checking instruction code %v == %v\n", req.InstructionCode, instCode)
		if req.InstructionCode == instCode {
			request = req
			err = nil
			continue
		}
	}
	if request == nil {
		return request, fmt.Errorf("invalid data instruction code %x", instCode)
	}
	return request, err
}

func (c *Client) send(req *Request) {
	for {
		if c.Sendable {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	respc := req.Responsec
	resp := Response{}
	log.Printf("sending request %#v\n", req.Message)
	var data [8]byte
	copy(data[:], req.Message)
	pSend := controlcan.CanObj{
		ID:         c.DevID,
		SendType:   1,
		RemoteFlag: 0,
		ExternFlag: 0,
		DataLen:    8,
		Data:       data,
	}
	if err := controlcan.Transmit(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
		pSend,
		1,
	); err != nil {
		log.Println(err)
		resp.Error = err
		respc <- resp
		return
	}
	log.Printf("request sent: %v\n", pSend)
	return
}

//
// no ack, com: message, [], 0x01
// ack com: message, 0x00, 0x01
// no ack, no com: message, [], []
func (c *Client) Send(
	message []byte,
	recExpected []byte,
	recIndex int,
	comExpected []byte,
	comIndex int,
) ([]byte, error) {
	code, err := c.getInstructionCode()
	if err != nil {
		return []byte{}, err
	}
	defer c.releaseInstructionCode(code)
	message = append(message, code)
	req := Request{
		InstructionCode: code,
		Message:         message,
		RecExpected:     recExpected,
		ComExpected:     comExpected,
		Responsec:       make(chan Response),
	}
	c.RequestQueue.Push(&req)
	if len(recExpected) > 0 {
		resp := <-req.Responsec
		status := resp.Message[recIndex]
		switch status {
		case STATUS_CODE_RECEIVED:
			log.Println("request received:", message)
		default:
			return resp.Message,
				fmt.Errorf("invalid status code %#v", status)
		}
	}
	resp := <-req.Responsec
	if len(comExpected) > 0 {
		status := resp.Message[comIndex]
		switch status {
		case STATUS_CODE_COMPLETED:
			return resp.Message, nil
		case STATUS_CODE_ERROR:
			return resp.Message,
				fmt.Errorf("unknown error when execute %#v", message)
		default:
			return resp.Message,
				fmt.Errorf("invalid status code %#v", status)
		}
	}
	c.ReceptionMap.Del(hex.EncodeToString([]byte{req.InstructionCode}))
	return resp.Message, resp.Error
}

func (c *Client) getInstructionCode() (code byte, err error) {
	var origin, update interface{}
	origin = false
	update = true
	for {
		log.Printf("allocating instruction code...")
		key, err := c.InstructionCodeMap.Replace(origin, update)
		if err == nil {
			codeSlice, err := hex.DecodeString(key)
			if err != nil {
				log.Println(err)
				return code, err
			}
			code = codeSlice[0]
			return code, nil
		}
		log.Printf("not enough instruction code, wait 5 seconds\n")
		time.Sleep(1000 * time.Millisecond)
	}
	return code, fmt.Errorf("not valid instruction code")
}

func (c *Client) occupyInstructionCode(code byte) {
	c.InstructionCodeMap.Set(
		hex.EncodeToString([]byte{code}),
		true,
	)
}

func (c *Client) releaseInstructionCode(code byte) {
	c.InstructionCodeMap.Set(
		hex.EncodeToString([]byte{code}),
		false,
	)
	log.Println("release instruction code: ", code)
}
