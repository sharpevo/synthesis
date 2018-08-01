package serialport_test

import (
	"github.com/tarm/serial"
	"posam/protocol/serialport"
	"time"

	"testing"
)

type MockPort struct {
	serialport.Port
}

func (p *MockPort) OpenPort(name string, baud int) (*serial.Port, error) {
	return &serial.Port{}, nil
}

func TestInstanceOperationOnMap(t *testing.T) {
	sp1 := &serialport.SerialPort{
		&MockPort{
			Port: serialport.Port{},
		},
		"/dev/ttyUSB0",
		9601,
		0x01,
		8,
		1,
		-1,
	}
	instance1 := sp1.Instance()

	sp2 := &serialport.SerialPort{
		&MockPort{
			Port: serialport.Port{},
		},
		"/dev/ttyUSB1",
		9602,
		0x01,
		8,
		1,
		-1,
	}
	instance2 := sp2.Instance()

	sp3 := &serialport.SerialPort{
		&MockPort{
			Port: serialport.Port{},
		},
		"/dev/ttyUSB0",
		9603,
		0x01,
		8,
		1,
		-1,
	}
	instance3 := sp3.Instance()
	t.Logf(
		"%p : %p : %p\n",
		instance1,
		instance2,
		instance3,
	)
	if instance1 == instance2 {
		t.Errorf(
			"Same serial port for the different devices\nEXPECT: %p\nGET: %p\n\n",
			instance1,
			instance2,
		)
	}
	if instance1 != instance3 {
		t.Errorf(
			"Different serial port for the same devices\nEXPECT: %p\nGET: %p\n\n",
			instance1,
			instance3,
		)
	}
}

func TestInstanceMapConcurrency(t *testing.T) {
	sp := &serialport.SerialPort{
		&MockPort{
			Port: serialport.Port{},
		},
		"/dev/ttyUSB0",
		9603,
		0x01,
		8,
		1,
		-1,
	}

	go func() {
		for {
			_ = sp.Instance()
		}
	}()

	go func() {
		for {
			openedPort := &serial.Port{}
			serialport.AddInstance(sp, openedPort)
		}
	}()

	select {
	case <-time.After(3 * time.Second):
	}
}
