package usbcan_test

import (
	"fmt"
	"posam/protocol/usbcan"
	"reflect"
	"testing"
)

func TestInstance(t *testing.T) { // {{{
	client1 := &usbcan.Client{
		Channel: nil,
		DevID:   1,
	}
	usbcan.ClientMap.Set("test", client1)

	cases := []struct {
		key    string
		client *usbcan.Client
	}{
		{
			"",
			client1, // return arbitrary items
		},
		{
			"test",
			client1,
		},
		{
			"not exist",
			nil,
		},
	}
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			actual := usbcan.Instance(c.key)
			expect := c.client
			if !reflect.DeepEqual(actual, expect) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestAddInstance(t *testing.T) { // {{{
	client1 := &usbcan.Client{
		Channel: nil,
		DevID:   1,
	}
	cases := []struct {
		client  *usbcan.Client
		existed bool
	}{
		{
			client1,
			false,
		},
		{
			client1,
			true,
		},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			_, ok := usbcan.AddInstance(c.client)
			if ok != c.existed {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.existed,
					ok,
				)
			}
		})
	}
} // }}}

func TestNewDevice(t *testing.T) { // {{{
	cases := []struct {
		devtype  int
		devindex int
		found    bool
	}{
		{
			0, 0, false,
		},
		{
			0, 0, true,
		},
		{
			1, 1, true,
		},
	}
	called := false
	originOpenDevice := usbcan.OpenDevice
	defer func() { usbcan.OpenDevice = originOpenDevice }()
	usbcan.OpenDevice = func(device *usbcan.Device) error {
		called = true
		if device.DeviceKey() == "1-1" {
			return fmt.Errorf("error expected")
		}
		t.Log("device opened", device.DeviceKey())
		return nil
	}

	c := cases[0]
	device, err := usbcan.NewDevice(c.devtype, c.devindex)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := usbcan.DeviceMap.Get(device.DeviceKey())
	if !ok || !called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("ok: %v; called: %v", true, true),
			fmt.Sprintf("ok: %v; called: %v", ok, called),
		)
	}
	called = false

	c = cases[1]
	device, err = usbcan.NewDevice(c.devtype, c.devindex)
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("called: %v", false),
			fmt.Sprintf("called: %v", called),
		)
	}
	called = false

	c = cases[2]
	device, err = usbcan.NewDevice(c.devtype, c.devindex)
	if err == nil {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"error",
			"nil error",
		)
	}
	_, ok = usbcan.DeviceMap.Get(device.DeviceKey())
	if ok || !called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("ok: %v; called: %v", false, true),
			fmt.Sprintf("ok: %v; called: %v", ok, called),
		)
	}
} // }}}

// maps pointer will be manipulated
// tests may failed if put it before others
func TestResetInstance(t *testing.T) { // {{{
	clientKey := "test reset instance"
	client1 := &usbcan.Client{
		Channel: nil,
		DevID:   3,
	}
	usbcan.ClientMap.Set(clientKey, client1)
	// never use usbcan.Client again
	// since it's the pointer to the *usbcan.clientMap,
	// not the  pointer to the usbcan.clientMap,
	client := usbcan.Instance(clientKey)
	if client == nil {
		t.Error("failed to add instance")
	}
	usbcan.ResetInstance()
	client = usbcan.Instance(clientKey)
	if client != nil {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			nil,
			client,
		)
	}
} // }}}

func TestDeviceKey(t *testing.T) { // {{{
	cases := []struct {
		devtype  int
		devindex int
		expected string
	}{
		{
			0, 0, "0-0",
		},
	}
	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			d := usbcan.Device{c.devtype, c.devindex}
			if d.DeviceKey() != c.expected {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expected,
					d.DeviceKey(),
				)
			}
		})
	}
} // }}}
