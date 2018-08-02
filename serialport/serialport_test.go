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
		Porter: &MockPort{
			Port: serialport.Port{},
		},
		Name:     "/dev/ttyUSB0",
		BaudRate: 9601,
		DataBits: 8,
		StopBits: 1,
		Parity:   -1,
	}
	instance1 := sp1.Instance()

	sp2 := &serialport.SerialPort{
		Porter: &MockPort{
			Port: serialport.Port{},
		},
		Name:     "/dev/ttyUSB1",
		BaudRate: 9602,
		DataBits: 8,
		StopBits: 1,
		Parity:   -1,
	}
	instance2 := sp2.Instance()

	sp3 := &serialport.SerialPort{
		Porter: &MockPort{
			Port: serialport.Port{},
		},
		Name:     "/dev/ttyUSB0",
		BaudRate: 9603,
		DataBits: 8,
		StopBits: 1,
		Parity:   -1,
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
		Porter: &MockPort{
			Port: serialport.Port{},
		},
		Name:     "/dev/ttyUSB0",
		BaudRate: 9603,
		DataBits: 8,
		StopBits: 1,
		Parity:   -1,
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
