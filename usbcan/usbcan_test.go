package usbcan_test

import (
	"posam/protocol/usbcan"
	"reflect"
	"testing"
)

func TestInstance(t *testing.T) {
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
}
