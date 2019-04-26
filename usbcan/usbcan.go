package usbcan

import (
	"controlcan"
	"encoding/hex"
	"fmt"
	"log"
	"posam/config"
	"posam/gui/uiutil"
	"posam/util"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"sync"
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

var (
	SEND_TIMEOUT          time.Duration
	WARN_TIMEOUT          time.Duration
	INTERRUPT_WHEN_WARN   = config.GetBool("can.transmission.interruptwhenwarning")
	RESEND_ALL            = config.GetBool("can.resend.all")
	RESEND_ONCE           = config.GetBool("can.resend.once")
	NOTIFY_RESEND_SUCCESS = config.GetBool("can.resend.notify.success")
	NOTIFY_RESEND_FAILURE = config.GetBool("can.resend.notify.failure")
)

var (
	deviceMap  = concurrentmap.NewConcurrentMap()
	channelMap = concurrentmap.NewConcurrentMap()
	clientMap  = concurrentmap.NewConcurrentMap()
)

func init() {
	config.SetDefault("can.transmission.timeout", 500)
	SEND_TIMEOUT = time.Duration(config.GetInt("can.transmission.timeout")) * time.Millisecond
	config.SetDefault("can.transmission.warningtimeout", 5000)
	WARN_TIMEOUT = time.Duration(config.GetInt("can.transmission.warningtimeout")) * time.Millisecond
}

func instance(key string) (client *Client) {
	if key == "" {
		for item := range clientMap.Iter() {
			if client == nil {
				client = item.Value.(*Client)
			}
		}
		return client
	} else {
		if clienti, ok := clientMap.Get(key); ok {
			return clienti.(*Client)
		}
	}
	return client
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
		client.Reset()
	}
	clientMap = concurrentmap.NewConcurrentMap()
	deviceMap = concurrentmap.NewConcurrentMap()
	channelMap = concurrentmap.NewConcurrentMap()
}

// A device is the abstraction of CANalyst device, which contains multiple
// channels. Note that the value of `DevIndex` is different between Liunx and
// Windows.
type Device struct {
	DevType  int
	DevIndex int
}

var newDevice = func(
	devType int,
	devIndex int,
) (*Device, error) {
	device := &Device{
		DevType:  devType,
		DevIndex: devIndex,
	}
	if d, found := deviceMap.Get(device.deviceKey()); found {
		return d.(*Device), nil
	}
	if err := OpenDevice(device); err != nil {
		return device, err
	}
	deviceMap.Set(device.deviceKey(), device)
	return device, nil

}

var OpenDevice = func(device *Device) error {
	return device.Open()
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
	return nil
}

func (d *Device) deviceKey() string {
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

	sendableLock sync.Mutex
	sendable     bool

	//receiveo sync.Once
	//sendo    sync.Once
	senderLaunched   bool
	receiverLaunched bool
}

func newChannel(
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
		sendable: true,
	}
	channel.DevType = devType
	channel.DevIndex = devIndex
	if _, err := newDevice(channel.DevType, channel.DevIndex); err != nil {
		return channel, err
	}
	if c, found := channelMap.Get(channel.ChannelKey()); found {
		return c.(*Channel), nil
	}
	channel.init()
	if err := StartChannel(channel); err != nil {
		return channel, err
	}
	channelMap.Set(channel.ChannelKey(), channel)
	return channel, nil
}

var StartChannel = func(channel *Channel) error {
	return channel.Start()
}

func (c *Channel) init() {
	c.RequestQueue = blockingqueue.NewBlockingQueue()
	c.ReceptionMap = concurrentmap.NewConcurrentMap()
	c.InstructionCodeMap = concurrentmap.NewConcurrentMap()
	c.loadInstructionCode()
}

func (c *Channel) loadInstructionCode() {
	for index := range [256]byte{} {
		c.InstructionCodeMap.Set(
			hex.EncodeToString([]byte{byte(index)}),
			false,
		)
	}
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
	util.Go(c.send)
	util.Go(c.receive)
	return nil
}

func (c *Channel) ChannelKey() string {
	return fmt.Sprintf("%v-%v-%v", c.DevType, c.DevIndex, c.CanIndex)
}

func (c *Channel) send() {
	for {
		reqi, err := c.RequestQueue.Pop()
		if err != nil {
			// TODO: error handling, e.g., insert to the response
			log.Printf("canalyst client sender %v terminated\n", c.deviceKey())
			return
		}
		c.ReceptionMap.Lock()
		req := reqi.(*Request)
		c.ReceptionMap.SetLockless(
			hex.EncodeToString([]byte{req.InstructionCode}),
			req,
		)
		TransmitRequest(c, req)
		c.ReceptionMap.Unlock()
		<-time.After(7 * time.Millisecond)
	}
}

var TransmitRequest = func(c *Channel, req *Request) {
	c.transmit(req)
}

func (c *Channel) transmit(req *Request) {
	c.untilSendable()
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
	pSendList := []controlcan.CanObj{pSend}
	if err := controlcan.Transmit(
		c.DevType,
		c.DevIndex,
		c.CanIndex,
		pSendList,
		len(pSendList),
	); err != nil {
		log.Println(err)
		resp.Error = err
		respc <- resp
		return
	}
	log.Printf("request sent: %v\n", pSend)
	req.TimeSent = time.Now()
	return
}

// }}}

// Transmit{{{
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
	defer c.ReceptionMap.Del(hex.EncodeToString([]byte{req.InstructionCode}))
	c.RequestQueue.Push(&req)
	if len(recExpected) > 0 {
		resp := <-req.Responsec
		if resp.Error != nil {
			return resp.Message, resp.Error
		}
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
		if resp.Error != nil {
			return resp.Message, resp.Error
		}
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
	return resp.Message, resp.Error
} // }}}

func (c *Channel) receive() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for _ = range ticker.C {
		pReceive, count, err := controlcan.Receive(
			c.DevType,
			c.DevIndex,
			c.CanIndex,
			100,
		)
		if err != nil || count == controlcan.MINUS_ONE {
			log.Printf("canalyst client receiver %v terminated\n", c.deviceKey())
			return
		}
		parseCanObjects(c, pReceive[:count])
		// TODO: not the best time to do that
		c.TryResend()
	}
}

var parseCanObjects = func(c *Channel, pReceive []controlcan.CanObj) {
	for _, canObj := range pReceive {
		data := make([]byte, len(canObj.Data))
		copy(data, canObj.Data[:])
		resp := Response{}
		req, err := findRequestByResponse(c, data, canObj.ID)
		if err != nil {
			log.Println(err)
			// TODO: notification
			continue
		}
		if NOTIFY_RESEND_SUCCESS {
			if req.ResendCount > 0 {
				msg := fmt.Sprintf(
					"request resent\nframe id: %v\ndata: %v\n",
					req.FrameId,
					req.Message,
				)
				uiutil.App.ShowMessageSlot(msg)
			}
		}
		resp.Message = data
		go func() {
			req.Responsec <- resp
		}()
	}
}

func (c *Channel) Reset() {
	if c == nil {
		return
	}
	c.RequestQueue.Reset()
	c.ReceptionMap = concurrentmap.NewConcurrentMap()
	c.loadInstructionCode()
}

func (c *Channel) Sendable() bool {
	c.sendableLock.Lock()
	defer c.sendableLock.Unlock()
	return c.sendable
}

func (c *Channel) SetSendable(sendable bool) {
	c.sendableLock.Lock()
	defer c.sendableLock.Unlock()
	c.sendable = sendable
}

func (c *Channel) untilSendable() {
	if !c.Sendable() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for _ = range ticker.C {
			if c.Sendable() {
				ticker.Stop()
				break
			}
		}
	}
}

// Helpers{{{

var parseFunctionCode = func(c *Channel, data []byte) (byte, error) {
	code := data[0]
	switch code {
	case 0xE0:
		c.SetSendable(false)
		return code, fmt.Errorf("mailbox is full (E0): %#v\n", data)
	case 0xE1:
		return code, fmt.Errorf("mailbox is overflow (E1): %#v\n", data)
	default:
	}
	return code, nil
}

var parseInstructionCode = func(c *Channel, data []byte) (code byte, err error) {
	code = data[7]
	if _, ok := c.InstructionCodeMap.Get(
		hex.EncodeToString([]byte{code}),
	); !ok {
		return code, fmt.Errorf("invalid instruction code")
	}
	return code, nil
}

var findRequestByResponse = func(
	c *Channel,
	data []byte,
	frameId int,
) (request *Request, err error) {
	_, err = parseFunctionCode(c, data[:])
	if err != nil {
		return request, err
	}
	instCode, err := parseInstructionCode(c, data)
	if err != nil {
		return request, err
	}
	for item := range c.ReceptionMap.Iter() {
		fmt.Println("findRequestByResponse: iter receptionmap")
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

func (c *Channel) TryResend() {
	now := time.Now()
	for item := range c.ReceptionMap.Iter() {
		reqi := item.Value
		req, ok := reqi.(*Request)
		if !ok {
			err := fmt.Errorf("invalid request: %#v", reqi)
			log.Println(err)
			continue
		}
		if INTERRUPT_WHEN_WARN {
			warningTimeout := req.TimeSent.Add(WARN_TIMEOUT)
			if warningTimeout.Before(now) {
				resp := Response{}
				resp.Error = fmt.Errorf("Warning: no response for request\nframe id: %v\ndata: %v\ntime: %v",
					req.FrameId,
					req.Message,
					req.TimeSent.Format("15:04:05.999999"),
				)
				go func() {
					req.Responsec <- resp
				}()
			}
		}
		if !RESEND_ALL {
			if req.Message[0] != 10 {
				continue
			}
		}
		timeout := req.TimeSent.Add(SEND_TIMEOUT)
		if timeout.Before(now) {
			c.resend(now, req)
		}
	}
}

func (c *Channel) resend(now time.Time, req *Request) {
	fmt.Printf("\n\n--------------------------------------------------\n\n")
	log.Printf(
		"error: can comm timeout\nframe id: %v\ncode: %v\ndata: %v\nresend count: %v\n",
		req.FrameId,
		req.InstructionCode,
		req.Message,
		req.ResendCount,
	)
	if RESEND_ONCE && req.ResendCount > 1 {
		// blocked Del
		c.ReceptionMap.DelLockless(hex.EncodeToString([]byte{req.InstructionCode}))
		resp := Response{}
		errmsg := fmt.Errorf(
			"error: failed to resend request\nframe id: %v\ndata: %v\n",
			req.FrameId,
			req.Message,
		)
		log.Println(errmsg)

		//resp.Error = errmsg
		resp.Message = []byte{0x0}
		go func() {
			req.Responsec <- resp
			if NOTIFY_RESEND_FAILURE {
				uiutil.App.ShowMessageSlot(errmsg.Error())
			}
		}()
		return
	}
	log.Printf(
		"resending...\nframe id: %v\ncode: %v\ndata: %v\n",
		req.FrameId,
		req.InstructionCode,
		req.Message,
	)
	req.ResendCount++
	req.TimeSent = now
	c.RequestQueue.Push(req)
	fmt.Printf("--------------------------------------------------\n\n")
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
	channel, err := newChannel(
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
	TimeSent        time.Time
	ResendCount     int
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

func (c *Client) Reset() {
	c.Channel.Reset()
}
