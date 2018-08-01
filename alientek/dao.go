package alientek

import (
	"fmt"
	"log"
	"posam/dao"
	"posam/protocol/serialport"
	"sync"
)

var deviceMap sync.Map

type Dao struct {
	DeviceAddress byte
	SerialPort    serialport.SerialPorter
}

func AddInstance(dao *Dao) {
	deviceMap.Store(string(dao.DeviceAddress), dao)
}

func Instance(address string) *Dao {
	if device, ok := deviceMap.Load(address); ok {
		return device.(*Dao)
	} else {
		return nil
	}
}

func (d *Dao) TurnOnLed() (resp string, err error) {
	req := d.makeRequest(LedOnUnit.Request())
	err = d.SerialPort.Send(req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}

	resp, err = d.checkResponse("Turn on LED", "sent", LedOnUnit.RecResp())
	if err != nil {
		return
	}
	resp, err = d.checkResponse("Turn on LED", "complete", LedOnUnit.ComResp())
	if err != nil {
		return
	}

	return
}

func (d *Dao) TurnOffLed() (resp string, err error) {
	req := d.makeRequest(LedOffUnit.Request())
	err = d.SerialPort.Send(req.Bytes())
	if err != nil {
		log.Println(err)
	}

	resp, err = d.checkResponse("Turn off LED", "sent", LedOffUnit.RecResp())
	if err != nil {
		return
	}
	resp, err = d.checkResponse("Turn off LED", "complete", LedOffUnit.ComResp())
	if err != nil {
		return
	}
	return
}

func (d *Dao) makeRequest(req dao.Request) dao.Request {
	req.Address = d.DeviceAddress
	return req
}

func (d *Dao) checkResponse(title string, action string, expect []byte) (resp string, err error) {
	log.Printf("%s: check %s response %x", title, action, expect)
	respBytes, err := d.SerialPort.Receive(expect)
	resp = fmt.Sprintf("%x", respBytes)
	if err != nil {
		return resp, fmt.Errorf("%s: failed to %s | %s", title, action, err)
	}
	return resp, nil
}
