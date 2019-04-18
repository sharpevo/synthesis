package usbcan_test

import (
	"encoding/hex"
	"fmt"
	"posam/protocol/usbcan"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewChannel(t *testing.T) { // {{{
	cases := []struct {
		devtype  int
		devindex int
		canindex int
		acccode  int
		accmask  int
		filter   int
		timing0  int
		timing1  int
		mode     int

		expectNewDeviceError string
		expectFound          bool
		expectStartError     string
	}{
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			"", false, "",
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			"", true, "",
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			"", true, "",
		},
	}
	dIndex := 0
	originNewDevice := usbcan.NewDevice
	defer func() { usbcan.NewDevice = originNewDevice }()
	usbcan.NewDevice = func(int, int) (*usbcan.Device, error) {
		defer func() { dIndex++ }()
		t.Log("dIndex", dIndex)
		msg := cases[dIndex].expectNewDeviceError
		if msg != "" {
			return nil, fmt.Errorf(msg)
		}
		return nil, nil
	}
	sIndex := 0
	originStartChannel := usbcan.StartChannel
	defer func() { usbcan.StartChannel = originStartChannel }()
	usbcan.StartChannel = func(*usbcan.Channel) error {
		defer func() { sIndex++ }()
		msg := cases[sIndex].expectStartError
		if msg != "" {
			return fmt.Errorf(msg)
		}
		return nil
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			if _, found := usbcan.ChannelMap.Get(fmt.Sprintf(
				"%v-%v-%v",
				c.devtype,
				c.devindex,
				c.canindex)); found != c.expectFound {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expectFound,
					found,
				)
			}
			_, err := usbcan.NewChannel(
				c.devtype, c.devindex, c.canindex, c.acccode, c.accmask,
				c.filter, c.timing0, c.timing1, c.mode,
			)
			if c.expectNewDeviceError != "" && err == nil {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expectNewDeviceError,
					err,
				)
			}
			if c.expectStartError != "" && err == nil {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expectStartError,
					err,
				)
			}
		})
	}
} // }}}

func TestChannelSend(t *testing.T) { // {{{
	cases := []struct {
		id  int
		req *usbcan.Request
		//response string
	}{
		{
			0,
			&usbcan.Request{InstructionCode: byte(5)},
		},
	}
	originTransmitRequest := usbcan.TransmitRequest
	defer func() { usbcan.TransmitRequest = originTransmitRequest }()
	usbcan.TransmitRequest = func(c *usbcan.Channel, req *usbcan.Request) {
		fmt.Println("fake transmit")
		return
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.id), func(t *testing.T) {
			channel := usbcan.Channel{}
			channel.RequestQueue = blockingqueue.NewBlockingQueue()
			channel.ReceptionMap = concurrentmap.NewConcurrentMap()
			go channel.Send() // receiver should be ahead of sender
			channel.RequestQueue.Push(c.req)
			<-time.After(1 * time.Second) // wait for sending
			channel.RequestQueue.Reset()
			reqReception, ok := channel.ReceptionMap.Get(
				hex.EncodeToString([]byte{c.req.InstructionCode}))
			fmt.Println(reqReception, ok)
			if !ok || !reflect.DeepEqual(reqReception, c.req) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.req,
					reqReception,
				)
			}
		})
	}
} // }}}

func TestUntilSendable(t *testing.T) { // {{{
	cases := []struct {
		leastTime int
	}{
		{
			1,
		},
		{
			2,
		},
	}
	for _, c := range cases {
		leastTime := time.Duration(c.leastTime)
		channel := &usbcan.Channel{}
		channel.SetSendable(false)
		start := time.Now()
		go func() {
			<-time.After(leastTime * time.Second)
			channel.SetSendable(true)
		}()
		channel.UntilSendable()
		end := time.Now()
		actual := end.Sub(start)
		if actual < leastTime {
			t.Errorf(
				"\nEXPECT: %v\n GET: %v\n\n",
				leastTime,
				actual,
			)
		}
	}
} // }}}

func TestChannelTransmit(t *testing.T) {
	cases := []struct {
		frameid     int
		recexpected []byte
		recindex    int
		comexpected []byte
		comindex    int

		errmsg      string
		ackResponse *usbcan.Response // pointer can be compared with nil
		comResponse *usbcan.Response // Message is sent as the final resp
	}{
		{ // no ack, no com: e.g., switch, humiture
			frameid:     1,
			recexpected: []byte{},
			recindex:    0,
			comexpected: []byte{},
			comindex:    0,
			errmsg:      "",
			comResponse: &usbcan.Response{
				Message: []byte{1},
			},
		},
		{ // ack com: e.g., switch advanced, motor
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED}, // 0x00
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED}, // 0x01
			comindex:    6,
			errmsg:      "",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_COMPLETED, 9},
			},
		},
		{ // no ack, com: e.g., rom read & write
			frameid:     1,
			recexpected: []byte{},
			recindex:    0,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    2,
			errmsg:      "",
			comResponse: &usbcan.Response{
				Message: []byte{0, 0, usbcan.STATUS_CODE_COMPLETED, 0, 0, 0, 0, 0},
			},
		},

		{ // ack com: ack error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{}, // stop sending extra message to channel
			comindex:    6,
			errmsg:      "ack error",
			ackResponse: &usbcan.Response{
				Error: fmt.Errorf("ack error"),
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
			},
		},
		{ // ack com: ack status error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{}, // stop sending extra message to channel
			comindex:    6,
			errmsg:      "invalid status code 0xe",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, 0xE, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
			},
		},
		{ // ack com: com error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "com error",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
				Error:   fmt.Errorf("com error"),
			},
		},
		{ // ack com: com error status code
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "unknown error when execute []byte{0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x2, 0x9", // omit instruction code
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_ERROR, 9},
			},
		},
		{ // ack com: com invalid status code
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "invalid status code 0xe", // omit instruction code
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, 0xE, 9},
			},
		},
	}
	next := make(chan struct{})
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			channel := &usbcan.Channel{}
			channel.Init()
			go func() {
				resp, err := channel.Transmit(
					c.frameid,
					c.comResponse.Message,
					c.recexpected,
					c.recindex,
					c.comexpected,
					c.comindex,
				)
				if err != nil && !strings.HasPrefix(err.Error(), c.errmsg) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.errmsg,
						err.Error(),
					)
				}
				if err == nil && !reflect.DeepEqual(resp, c.comResponse.Message) { // not same when error
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.comResponse.Message,
						resp,
					)
				}
				next <- struct{}{}
			}()
			reqi, _ := channel.RequestQueue.Pop()
			req := reqi.(*usbcan.Request)
			if c.ackResponse != nil {
				req.Responsec <- *c.ackResponse
			}
			if c.errmsg == "" || // not sending com when error occured, especially ack error
				len(c.comexpected) > 0 { // still sending when com error expected
				req.Responsec <- *c.comResponse
			}
			<-next
		})
	}
}
