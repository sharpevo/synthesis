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
	originOpenDevice := OpenDevice
	defer func() { OpenDevice = originOpenDevice }()
	OpenDevice = func(device *Device) error {
		called = true
		if device.DeviceKey() == "1-1" {
			return fmt.Errorf("error expected")
		}
		t.Log("device opened", device.DeviceKey())
		return nil
	}

	c := cases[0]
	device, err := newDevice(c.devtype, c.devindex)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := DeviceMap.Get(device.DeviceKey())
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
	_, ok = DeviceMap.Get(device.DeviceKey())
	if ok || !called {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			fmt.Sprintf("ok: %v; called: %v", false, true),
			fmt.Sprintf("ok: %v; called: %v", ok, called),
		)
	}
} // }}}
