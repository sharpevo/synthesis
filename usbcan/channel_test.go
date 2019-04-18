package usbcan_test

import (
	"encoding/hex"
	"fmt"
	"posam/protocol/usbcan"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
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

func TestChannelSend(t *testing.T) {
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
}

func TestUntilSendable(t *testing.T) {
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
}
