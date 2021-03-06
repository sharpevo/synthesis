package serialport

import (
	"bytes"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"synthesis/internal/protocol/modbus"
	"synthesis/pkg/concurrentmap"
)

type Porter interface {
	OpenPort(string, int) (*serial.Port, error)
}

type SerialPorter interface {
	Instance() *serial.Port
	Send([]byte) error
	Receive([]byte) ([]byte, error)
}

type Port struct {
}

func (p *Port) OpenPort(name string, baud int) (*serial.Port, error) {
	log.Printf("Opening serial port %q...", name)
	c := &serial.Config{
		Name: name,
		Baud: baud,
	}
	openedPort, err := serial.OpenPort(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return openedPort, nil
}

var instanceMap *concurrentmap.ConcurrentMap

func init() {
	instanceMap = concurrentmap.NewConcurrentMap()
}

type SerialPort struct {
	Porter
	Name     string
	BaudRate int

	DeviceAddress byte

	DataBits int
	StopBits int
	Parity   int
}

// TODO: thread safe
// TODO: return busy error
func (s *SerialPort) Instance() *serial.Port {
	if instance, ok := instanceMap.Get(s.Name); ok {
		return instance.(*serial.Port)
	}
	if s.Porter == nil {
		s.Porter = &Port{}
	}
	p, _ := s.OpenPort(s.Name, s.BaudRate)
	return s.addInstance(p)
}

func (s *SerialPort) addInstance(openedPort *serial.Port) *serial.Port {
	instance := instanceMap.Set(s.Name, openedPort)
	return instance.(*serial.Port)
}

// data: with address of device
func (s *SerialPort) Send(data []byte) (err error) {
	s.Instance().Flush()
	modbus.AppendCRC(&data)
	if _, err = s.Instance().Write(data); err != nil {
		return err
	}
	return nil
}

func (s *SerialPort) Receive(expect []byte) (resp []byte, err error) {
	max := len(expect)
	buf := make([]byte, max)
	cnt := 0
	for {
		n, err := s.Instance().Read(buf)
		if err != nil {
			return resp, err
		}

		cnt += n
		resp = append(resp, buf[:n]...)
		if cnt >= max || n == 0 {
			break
		}
	}

	log.Printf("%x | %x", expect, resp)
	if !bytes.Equal(expect, resp) {
		s.Instance().Flush()
		return resp, fmt.Errorf(
			"invalid response code %x (%x)",
			resp,
			expect,
		)
	}
	return

}

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
