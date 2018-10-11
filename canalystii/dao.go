package canalystii

import (
	"fmt"
	"log"
	"posam/interpreter"
	"posam/protocol/usbcan"
	"posam/util/concurrentmap"
	"strconv"
)

const (
	NAME         = "CANALYSTII"
	DEVICE_TYPE  = "DEVICE_TYPE"
	DEVICE_INDEX = "DEVICE_INDEX"
	FRAME_ID     = "FRAME_ID"
	CAN_INDEX    = "CAN_INDEX"

	ACC_CODE = "ACC_CODE"
	ACC_MASK = "ACC_MASK"
	FILTER   = "FILTER"
	TIMING0  = "TIMING0"
	TIMING1  = "TIMING1"
	MODE     = "MODE"

	IDNAME = FRAME_ID
)

var CONN_ATTRIBUTES = []string{
	DEVICE_TYPE,
	DEVICE_INDEX,
	FRAME_ID,
	CAN_INDEX,
	ACC_CODE,
	ACC_MASK,
	FILTER,
	TIMING0,
	TIMING1,
	MODE,
}

var InstructionMap interpreter.InstructionMapt
var deviceMap *concurrentmap.ConcurrentMap

func init() {
	InstructionMap = make(interpreter.InstructionMapt)
	ResetInstance()
}

type Dao struct {
	id           string
	UsbcanClient usbcan.Clienter
}

func NewDao(
	devType string,
	devIndex string,
	devID string,
	canIndex string,
	accCode string,
	accMask string,
	filter string,
	timing0 string,
	timing1 string,
	mode string,
) (dao *Dao, err error) {
	devTypeInt, err := strconv.Atoi(devType)
	if err != nil {
		return dao, err
	}
	devIndexInt, err := strconv.Atoi(devIndex)
	if err != nil {
		return dao, err
	}
	devIDInt, err := strconv.ParseUint(devID, 0, 32)
	if err != nil {
		return dao, err
	}
	canIndexInt, err := strconv.Atoi(canIndex)
	if err != nil {
		return dao, err
	}
	accCodeInt, err := strconv.ParseUint(accCode, 0, 32)
	if err != nil {
		return dao, err
	}
	accMaskInt, err := strconv.ParseUint(accMask, 0, 32)
	if err != nil {
		return dao, err
	}
	filterInt, err := strconv.Atoi(filter)
	if err != nil {
		return dao, err
	}
	timing0Int, err := strconv.ParseUint(timing0, 0, 32)
	if err != nil {
		return dao, err
	}
	timing1Int, err := strconv.ParseUint(timing1, 0, 32)
	if err != nil {
		return dao, err
	}
	modeInt, err := strconv.Atoi(mode)
	if err != nil {
		return dao, err
	}
	usbcanClient, err := usbcan.NewClient(
		devTypeInt,
		devIndexInt,
		int(devIDInt),
		canIndexInt,
		int(accCodeInt),
		int(accMaskInt),
		filterInt,
		int(timing0Int),
		int(timing1Int),
		modeInt,
	)
	if err != nil {
		return dao, err
	}
	dao = &Dao{
		UsbcanClient: usbcanClient,
	}
	if err = dao.SetID(fmt.Sprintf("%v", devIDInt)); err != nil {
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
	usbcan.ResetInstance()
}

func Instance(id string) (dao *Dao, err error) {
	if id == "" {
		for device := range deviceMap.Iter() {
			return device.Value.(*Dao), nil
		}
		return dao, fmt.Errorf("empty instance pool")
	}
	if device, ok := deviceMap.Get(id); ok {
		return device.(*Dao), nil
	}
	return dao, fmt.Errorf("invalid device id %q", id)
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

func (d *Dao) ReadOxygenConc() (resp interface{}, err error) {
	req := SensorOxygenConcUnit.Request()
	resp, err = d.UsbcanClient.Send(
		req.Bytes(),
		SensorOxygenConcUnit.RecResp(),
		SensorOxygenConcUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return resp, nil
}
