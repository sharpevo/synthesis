package alientek_test

import (
	"log"
	"reflect"
	"strings"
	"synthesis/dao/alientek"
	"synthesis/protocol/serialport"
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
	alientek.ResetInstance()
	alientekDao, err := alientek.NewDao(
		"/dev/ttyUSB0",
		9600,
		8,
		1,
		-1,
		0x01,
	)
	if err != nil {
		t.Fatal(err)
	}

	alientekDao.TurnOnLed()

}

func TestInstanceOperation(t *testing.T) {
	alientek.ResetInstance()
	alientekDao, err := alientek.NewDao(
		"/dev/ttyUSB0",
		9600,
		8,
		1,
		-1,
		0x01,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(
		alientekDao,
		alientek.Instance(alientekDao.ID()),
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
			_, err := alientek.NewDao(
				"/dev/ttyUSB0",
				9600,
				8,
				1,
				-1,
				0x01,
			)
			if !strings.Contains(err.Error(), "existed") {
				t.Errorf(err.Error())
			}
		}
	}()
	select {
	case <-time.After(3 * time.Second):
	}
}
