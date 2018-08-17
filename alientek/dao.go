package alientek

import (
	"fmt"
	"log"
	"posam/dao"
	"posam/protocol/serial"
	//"posam/protocol/serialport"
	"posam/util/concurrentmap"
)

var deviceMap *concurrentmap.ConcurrentMap

func init() {
	ResetInstance()
}

type Dao struct {
	id           string
	AddressByte  byte
	SerialClient serial.Clienter
}

func NewDao(
	name string,
	baud int,
	databits int,
	stopbits int,
	parity int,
	addressbyte byte,
) (*Dao, error) {
	serialClient, err := serial.NewClient(
		name,
		baud,
		databits,
		stopbits,
		parity,
	)
	if err != nil {
		return nil, err
	}
	dao := &Dao{
		AddressByte:  addressbyte,
		SerialClient: serialClient,
	}
	err = dao.SetID(string(addressbyte))
	if err != nil {
		return dao, err
	}
	AddInstance(dao)
	return dao, nil
}

func AddInstance(dao *Dao) {
	deviceMap.Set(dao.ID(), dao)
}

func ResetInstance() {
	deviceMap = concurrentmap.NewConcurrentMap()
	serial.ResetInstance()
}

func Instance(address string) *Dao {
	if address == "" {
		for device := range deviceMap.Iter() {
			return device.Value.(*Dao)
		}
	} else {
		if device, ok := deviceMap.Get(address); ok {
			return device.(*Dao)
		} else {
			return nil
		}
	}
	log.Println("nil device instance")
	return nil
}

func (d *Dao) ID() string {
	return d.id
}

func (d *Dao) SetID(id string) error {
	if _, ok := deviceMap.Get(id); ok {
		return fmt.Errorf("ID %q is duplicated", id)
	}
	d.id = id
	return nil
}

func (d *Dao) TurnOnLed() (resp interface{}, err error) {
	req := d.makeRequest(LedOnUnit.Request())
	resp, err = d.SerialClient.Send(
		req.Bytes(),
		LedOnUnit.RecResp(),
		LedOnUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return
}

func (d *Dao) TurnOffLed() (resp interface{}, err error) {
	req := d.makeRequest(LedOffUnit.Request())
	resp, err = d.SerialClient.Send(
		req.Bytes(),
		LedOffUnit.RecResp(),
		LedOffUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return
}

func (d *Dao) makeRequest(req dao.Request) dao.Request {
	req.Address = d.AddressByte
	return req
}
