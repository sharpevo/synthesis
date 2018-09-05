package alientek

import (
	"encoding/hex"
	"fmt"
	"log"
	"posam/dao"
	"posam/protocol/serial"
	"strconv"
	//"posam/protocol/serialport"
	"posam/interpreter"
	"posam/util/concurrentmap"
)

const (
	NAME           = "ALIENTEK"
	DEVICE_CODE    = "DEVICE_CODE"
	DEVICE_NAME    = "DEVICE_NAME"
	BAUD_RATE      = "BAUD_RATE"
	CHARACTER_BITS = "CHARACTER_BITS"
	STOP_BITS      = "STOP_BITS"
	PARITY         = "PARITY"
	IDNAME         = DEVICE_CODE
)

var CONN_ATTRIBUTES = []string{
	DEVICE_CODE,
	DEVICE_NAME,
	BAUD_RATE,
	CHARACTER_BITS,
	STOP_BITS,
	PARITY,
}

var InstructionMap interpreter.InstructionMapt

var deviceMap *concurrentmap.ConcurrentMap

func init() {
	InstructionMap = make(interpreter.InstructionMapt)
	ResetInstance()
}

type Dao struct {
	id           string
	AddressByte  byte
	SerialClient serial.Clienter
}

func NewDao(
	name string,
	baud string,
	character string,
	stop string,
	parity string,
	deviceCode string,
) (dao *Dao, err error) {

	baudInt, err := strconv.Atoi(baud)
	if err != nil {
		return dao, err
	}
	characterInt, err := strconv.Atoi(character)
	if err != nil {
		return dao, err
	}
	stopInt, err := strconv.Atoi(stop)
	if err != nil {
		return dao, err
	}
	parityInt, err := strconv.Atoi(parity)
	if err != nil {
		return dao, err
	}
	codeBytes, err := hex.DecodeString(deviceCode)
	if err != nil {
		return dao, err
	}

	serialClient, err := serial.NewClient(
		name,
		baudInt,
		characterInt,
		stopInt,
		parityInt,
	)
	if err != nil {
		return nil, err
	}
	dao = &Dao{
		AddressByte:  codeBytes[0],
		SerialClient: serialClient,
	}
	err = dao.SetID(deviceCode)
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
