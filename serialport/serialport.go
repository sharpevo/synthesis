package serialport

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"posam/protocol/modbus"
)

type SerialPorter interface {
	Send([]byte) error
	Receive(string) ([]byte, error)
}

var instanceMap map[string]*serial.Port

type SerialPort struct {
	Name     string
	BaudRate int

	DeviceAddress byte

	DataBits int
	StopBits int
	Parity   int
}

func init() {
	instanceMap = make(map[string]*serial.Port)
}

// TODO: thread safe
// TODO: return busy error
func (s *SerialPort) Instance() *serial.Port {
	if instanceMap[s.Name] != nil {
		return instanceMap[s.Name]
	}
	log.Printf("Opening serial port %q...", s.Name)
	c := &serial.Config{
		Name: s.Name,
		Baud: s.BaudRate,
	}
	p, err := serial.OpenPort(c)
	if err != nil {
		log.Println(err)
		return nil
	}
	instanceMap[s.Name] = p
	return instanceMap[s.Name]

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

func (s *SerialPort) Receive(expect string) (resp []byte, err error) {
	max := len(expect) / 2
	buf := make([]byte, max)
	cnt := 0
	for {
		n, err := s.Instance().Read(buf)
		if err != nil {
			return resp, err
		}

		resp = append(resp, buf[:n]...)
		if cnt >= max || n == 0 {
			break
		}
	}

	if expect != toHexString(resp) {
		return resp, fmt.Errorf(
			"invalid response cocde %s (%s)",
			resp,
			toHexString(resp),
		)
	}
	return

}

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
