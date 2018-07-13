package alientek

import (
	"log"
	"posam/dao"
	"posam/protocol/serialport"
)

type Dao struct {
	DeviceAddress byte
	SerialPort    serialport.SerialPorter
}

func (d *Dao) TurnOnLed() {
	req := d.makeRequest(LedOnUnit.Request())
	err := d.SerialPort.Send(req.Bytes())
	if err != nil {
		log.Println(err)
	}
}

func (d *Dao) TurnOffLed() {
	req := d.makeRequest(LedOffUnit.Request())
	err := d.SerialPort.Send(req.Bytes())
	if err != nil {
		log.Println(err)
	}
}

func (d *Dao) makeRequest(req dao.Request) dao.Request {
	req.Address = d.DeviceAddress
	return req
}
