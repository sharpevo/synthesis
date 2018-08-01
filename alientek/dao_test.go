package alientek_test

import (
	"log"
	"posam/dao/alientek"
	"posam/protocol/serialport"
	"reflect"
	"testing"
	"time"
)

type MockSerialPort struct {
	serialport.SerialPort
}

func (m *MockSerialPort) Send(data []byte) error {
	log.Println("sent")
	return nil
}

func (m *MockSerialPort) Receive(expect []byte) (data []byte, err error) {
	log.Println("received")
	return
}

func TestTurnOnLed(t *testing.T) {

	alientekDao := alientek.Dao{
		DeviceAddress: 0x01,
		SerialPort: &MockSerialPort{
			serialport.SerialPort{
				Name:     "/dev/ttyUSB0",
				BaudRate: 9600,
				DataBits: 8,
				StopBits: 1,
				Parity:   -1,
			},
		},
	}

	alientekDao.TurnOnLed()

}

func TestInstanceOperation(t *testing.T) {

	alientekDao := &alientek.Dao{
		DeviceAddress: 0x01,
		SerialPort: &MockSerialPort{
			serialport.SerialPort{
				Name:     "/dev/ttyUSB0",
				BaudRate: 9600,
				DataBits: 8,
				StopBits: 1,
				Parity:   -1,
			},
		},
	}
	alientek.AddInstance(alientekDao)
	if !reflect.DeepEqual(
		alientekDao,
		alientek.Instance(string(alientekDao.DeviceAddress)),
	) {
		t.Errorf("Failed to get instance from deviceMap")
	}
}

func TestInstanceConcurrency(t *testing.T) {
	go func() {
		for {
			alientek.Instance("02")
		}
	}()
	go func() {
		for {
			alientek.AddInstance(&alientek.Dao{
				DeviceAddress: 0x01,
				SerialPort: &MockSerialPort{
					serialport.SerialPort{
						Name:     "/dev/ttyUSB0",
						BaudRate: 9600,
						DataBits: 8,
						StopBits: 1,
						Parity:   -1,
					},
				},
			})
		}
	}()
	select {
	case <-time.After(3 * time.Second):
	}
}
