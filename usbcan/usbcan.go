package usbcan

import (
	"controlcan"
	"encoding/hex"
	"fmt"
	"log"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	//"sync"
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
var channelMap *concurrentmap.ConcurrentMap
var clientMap *concurrentmap.ConcurrentMap

var receptMap map[string]*Client
var connMap map[string]*Client

func init() {
	clientMap = concurrentmap.NewConcurrentMap()
	deviceMap = concurrentmap.NewConcurrentMap()
	channelMap = concurrentmap.NewConcurrentMap()

	receptMap = make(map[string]*Client)
	connMap = make(map[string]*Client)
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
		//client.RequestQueue.Reset()
	}
	clientMap = concurrentmap.NewConcurrentMap()
	deviceMap = concurrentmap.NewConcurrentMap()
	channelMap = concurrentmap.NewConcurrentMap()
	receptMap = make(map[string]*Client)
	connMap = make(map[string]*Client)
}

type Device struct {
	DevType  int
	DevIndex int
}

func NewDevice(
	devType int,
	devIndex int,
) (*Device, error) {
	device := &Device{
		DevType:  devType,
		DevIndex: devIndex,
	}
	if d, found := deviceMap.Get(device.DeviceKey()); found {
		return d.(*Device), nil
	}
	if err := device.Open(); err != nil {
		log.Println(err)
		return device, err
	}
	deviceMap.Set(device.DeviceKey(), device)
	return device, nil

}

func (d *Device) Open() error {
	log.Printf(
		"Opening device(type %v, index %#v)...\n",
		d.DevType,
		d.DevIndex,
	)
	if err := controlcan.OpenDevice(
		int(d.DevType),
		int(d.DevIndex),
		0,
	); err != nil {
		return err
	}
	//deviceMap.Set(d.DeviceKey(), true)
	return nil
}

func (d *Device) DeviceKey() string {
	return fmt.Sprintf("%v-%v", d.DevType, d.DevIndex)
}

type Channel struct {
	Device   //   *Device
	CanIndex int
	AccCode  int
	AccMask  int
	Filter   int
	Timing0  int
	Timing1  int
	Mode     int

	RequestQueue       *blockingqueue.BlockingQueue
	ReceptionMap       *concurrentmap.ConcurrentMap
	InstructionCodeMap *concurrentmap.ConcurrentMap

	Sendable bool

	//receiveo sync.Once
	//sendo    sync.Once
	senderLaunched   bool
	receiverLaunched bool
}

func NewChannel(
	devType int,
	devIndex int,
	canIndex int,
	accCode int,
	accMask int,
	filter int,
	timing0 int,
	timing1 int,
	mode int,
) (*Channel, error) {
	channel := &Channel{
		CanIndex: canIndex,
		AccCode:  accCode,
		AccMask:  accMask,
		Filter:   filter,
		Timing0:  timing0,
		Timing1:  timing1,
		Mode:     mode,
		Sendable: true,
	}
	channel.DevType = devType
	channel.DevIndex = devIndex
	_, err := NewDevice(channel.DevType, channel.DevIndex)
	if err != nil {
		log.Println(err)
		return channel, err
	}
	if c, found := channelMap.Get(channel.ChannelKey()); found {
		return c.(*Channel), nil
	}
	channel.RequestQueue = blockingqueue.NewBlockingQueue()
	channel.ReceptionMap = concurrentmap.NewConcurrentMap()
	channel.InstructionCodeMap = concurrentmap.NewConcurrentMap()
	for index := range [256]byte{} {
		channel.InstructionCodeMap.Set(
			hex.EncodeToString([]byte{byte(index)}),
			false,
		)
	}
	if err := channel.Start(); err != nil {
		log.Println(err)
		return channel, err
	}
	channelMap.Set(channel.ChannelKey(), channel)
	return channel, nil
}

func (c *Channel) Start() error {
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
	if err := controlcan.InitCAN(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
		config,
	); err != nil {
		return err
	}
	log.Printf("Starting CAN(type %v, index %v, can %v)...\n", c.DevType, c.DevIndex, c.CanIndex)
	if err := controlcan.StartCAN(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
	); err != nil {
		return err
	}
	go c.send()
	go c.receive()
	return nil
}

func (c *Channel) ChannelKey() string {
	return fmt.Sprintf("%v-%v-%v", c.DevType, c.DevIndex, c.CanIndex)
}

func (c *Channel) send() error { // {{{
	//func (c *Channel) send() {
	//c.sendo.Do(func() {
	//if c.senderLaunched {
	//return nil
	//}
	//c.senderLaunched = true
	for {
		reqi, err := c.RequestQueue.Pop()
		log.Println("processing", reqi)
		if err != nil {
			log.Printf("canalyst client sender %v terminated\n", c.DeviceKey())
			return err
		}
		req := reqi.(*Request)
		c.ReceptionMap.Set(
			hex.EncodeToString([]byte{req.InstructionCode}),
			req,
		)
		c.transmit(req)
	}
	//})
}

func (c *Channel) transmit(req *Request) {
	if !c.Sendable {
		ticker := time.NewTicker(100 * time.Millisecond)
		for _ = range ticker.C {
			if c.Sendable {
				ticker.Stop()
				break
			}
		}
	}
	respc := req.Responsec
	resp := Response{}
	log.Printf("sending request %#v\n", req.Message)
	var data [8]byte
	copy(data[:], req.Message)
	pSend := controlcan.CanObj{
		ID:         req.FrameId,
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
} // }}}

// Transmit{{{
// no ack, com: message, [], 0x01
// ack com: message, 0x00, 0x01
// no ack, no com: message, [], []
func (c *Channel) Transmit(
	frameId int,
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
		FrameId:         frameId,
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
} // }}}

func (c *Channel) receive() {
	//if c.receiverLaunched {
	//return
	//}
	//c.receiverLaunched = true
	//c.receiveo.Do(func() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for _ = range ticker.C {
		pReceive := make(
			[]controlcan.CanObj,
			controlcan.FRAME_LENGTH_OF_RECEPTION,
		)
		count, err := controlcan.Receive(
			c.DevType,
			c.DevIndex,
			c.CanIndex,
			pReceive,
			100,
		)
		if err != nil || count < 0 {
			log.Printf(
				"canalyst client receiver %v terminated\n",
				c.DeviceKey(),
			)
			return
		}
		if count == 0 {
			continue
		}
		log.Printf("data received: %#v\n", pReceive[:count])

		for _, canObj := range pReceive[:count] {
			//devId := string(canObj.ID)
			devId := canObj.ID
			data := make([]byte, len(canObj.Data))
			copy(data, canObj.Data[:])
			resp := Response{}
			req, err := c.findRequestByResponse(data, devId)
			if err != nil {
				log.Println(err)
				// TODO: notification
				continue
			}
			resp.Message = data
			req.Responsec <- resp
		}
	}
	//})
}

// Helpers{{{

func (c *Channel) parseFunctionCode(data []byte) (byte, error) {
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

func (c *Channel) parseInstructionCode(data []byte) (code byte, err error) {
	code = data[7]
	if _, ok := c.InstructionCodeMap.Get(
		hex.EncodeToString([]byte{code}),
	); !ok {
		return code, fmt.Errorf("invalid instruction code")
	}
	return code, nil
}

func (c *Channel) findRequestByResponse(data []byte, frameId int) (request *Request, err error) {
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
		log.Printf("checking frame id %v == %v\n", req.FrameId, frameId)
		if req.FrameId != frameId {
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

func (c *Channel) getInstructionCode() (code byte, err error) {
	var origin, update interface{}
	origin = false
	update = true
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
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()
	for _ = range ticker.C {
		log.Printf("waiting for instruction code...")
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
	}
	return code, fmt.Errorf("not valid instruction code")
}

func (c *Channel) releaseInstructionCode(code byte) {
	c.InstructionCodeMap.Set(
		hex.EncodeToString([]byte{code}),
		false,
	)
	log.Println("release instruction code: ", code)
} // }}}

type Clienter interface {
	//connect() error
	Send([]byte, []byte, int, []byte, int) ([]byte, error)
}

type Client struct {
	Channel *Channel
	DevID   int // used as frame id
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
	client := &Client{}
	channel, err := NewChannel(
		devType,
		devIndex,
		canIndex,
		accCode,
		accMask,
		filter,
		timing0,
		timing1,
		mode,
	)
	if err != nil {
		log.Println(err)
		return client, err
	}
	client.DevID = devID
	if c, found := channelMap.Get(
		fmt.Sprintf("%v-%v-%v", devType, devIndex, canIndex),
	); found {
		channel = c.(*Channel)
	}
	client.Channel = channel
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	//go client.Channel.receive()
	//go client.Channel.send()
	return client, nil
}

type Request struct {
	FrameId         int
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
	return c.Channel.Transmit(
		c.DevID,
		message,
		recExpected,
		recIndex,
		comExpected,
		comIndex,
	)
}
