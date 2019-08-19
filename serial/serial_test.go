package serial_test

import (
	"fmt"
	"posam/protocol/serial"
	"reflect"
	"testing"
)

func TestInstanceOperationOnMap(t *testing.T) {
	s1, err := serial.NewClient(
		"/dev/ttyUSB0",
		9601,
		8,
		1,
		-1,
	)
	if err != nil {
		fmt.Println(err)
	}
	s2, err := serial.NewClient(
		"/dev/ttyUSB1",
		9602,
		8,
		1,
		-1,
	)
	if err != nil {
		fmt.Println(err)
	}
	s3, err := serial.NewClient(
		"/dev/ttyUSB0",
		9603,
		8,
		1,
		-1,
	)
	if err != nil {
		fmt.Println(err)
	}

	if !reflect.DeepEqual(s1, s3) {
		t.Errorf(
			"Different serial client for the same devices\nEXPECT: %p\nGET: %p\n\n",
			s1,
			s2,
		)
	}
	if reflect.DeepEqual(s1, s2) {
		t.Errorf(
			"Same serial client for the different devices\nEXPECT: %p\nGET: %p\n\n",
			s1,
			s2,
		)
	}
}
