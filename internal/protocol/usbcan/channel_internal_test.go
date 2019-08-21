package usbcan

import (
	"encoding/hex"
	"fmt"
	"synthesis/util/blockingqueue"
	"synthesis/util/concurrentmap"
	"reflect"
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
	originNewDevice := newDevice
	defer func() { newDevice = originNewDevice }()
	newDevice = func(int, int) (*Device, error) {
		defer func() { dIndex++ }()
		t.Log("dIndex", dIndex)
		msg := cases[dIndex].expectNewDeviceError
		if msg != "" {
			return nil, fmt.Errorf(msg)
		}
		return nil, nil
	}
	sIndex := 0
	originStartChannel := startChannel
	defer func() { startChannel = originStartChannel }()
	startChannel = func(*Channel) error {
		defer func() { sIndex++ }()
		msg := cases[sIndex].expectStartError
		if msg != "" {
			return fmt.Errorf(msg)
		}
		return nil
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			if _, found := ChannelMap.Get(fmt.Sprintf(
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
			_, err := newChannel(
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
		req *Request
		//response string
	}{
		{
			0,
			&Request{InstructionCode: byte(5)},
		},
	}
	originTransmitRequest := transmitRequest
	defer func() { transmitRequest = originTransmitRequest }()
	transmitRequest = func(c *Channel, req *Request) {
		fmt.Println("fake transmit")
		return
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.id), func(t *testing.T) {
			channel := Channel{}
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
			200,
		},
		{
			300,
		},
	}
	for _, c := range cases {
		leastTime := time.Duration(c.leastTime)
		channel := &Channel{}
		channel.setSendable(false)
		start := time.Now()
		go func() {
			<-time.After(leastTime * time.Millisecond)
			channel.setSendable(true)
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
