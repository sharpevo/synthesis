package usbcan

import (
	"fmt"
	"testing"
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
