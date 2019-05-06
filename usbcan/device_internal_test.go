package usbcan

import (
	"fmt"
	//"log"
	"testing"
)

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
	originOpenDevice := openDevice
	defer func() { openDevice = originOpenDevice }()
	openDevice = func(device *Device) error {
		called = true
		if device.deviceKey() == "1-1" {
			return fmt.Errorf("error expected")
		}
		t.Log("device opened", device.deviceKey())
		return nil
	}

	c := cases[0]
	device, err := newDevice(c.devtype, c.devindex)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := DeviceMap.Get(device.deviceKey())
	if !ok || !called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("ok: %v; called: %v", true, true),
			fmt.Sprintf("ok: %v; called: %v", ok, called),
		)
	}
	called = false

	c = cases[1]
	device, err = newDevice(c.devtype, c.devindex)
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
	device, err = newDevice(c.devtype, c.devindex)
	if err == nil {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"error",
			"nil error",
		)
	}
	_, ok = DeviceMap.Get(device.deviceKey())
	if ok || !called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("ok: %v; called: %v", false, true),
			fmt.Sprintf("ok: %v; called: %v", ok, called),
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
			d := Device{c.devtype, c.devindex}
			if d.deviceKey() != c.expected {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expected,
					d.deviceKey(),
				)
			}
		})
	}
} // }}}
