package tml

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddInstance(t *testing.T) { // {{{
	cases := []struct {
		name  string
		found bool
	}{
		{
			"1",
			false,
		},
		{
			"2",
			false,
		},
		{
			"1",
			true,
		},
	}
	originLaunchClient := launchClient
	defer func() { launchClient = originLaunchClient }()
	launchClient = func(client *Client) {}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			client := &Client{Name: c.name}
			rclient, found := addInstance(client)
			if !reflect.DeepEqual(rclient, client) ||
				found != c.found {
				t.Errorf(
					"\nEXPECT: %v %v\n GET: %v %v\n\n",
					c.found, client,
					found, rclient,
				)
			}

		})
	}
} // }}}
