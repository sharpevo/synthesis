package alientek_test

import (
	"log"
	"posam/dao/alientek"
	"posam/protocol/serialport"
	"testing"
)

type MockSerialPort struct {
	serialport.SerialPort
}

func (m *MockSerialPort) Send(data []byte) error {
	log.Println("sent")
	return nil
}

func (m *MockSerialPort) Receive(expect string) (data []byte, err error) {
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
