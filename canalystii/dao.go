package canalystii

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
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

func (d *Dao) MoveRelative(
	motorCode string,
	speed string,
	position string,
) (resp interface{}, err error) {
	motorCodeBytes, err := stringToUint8Bytes(motorCode)
	if err != nil {
		return resp, err
	}
	speedBytes, _, err := stringToUint16Bytes(speed)
	if err != nil {
		return resp, err
	}
	posBytes, directionBytes, err := positionBytes(position)
	if err != nil {
		return resp, err
	}
	req := MotorMoveRelativeUnit.Request()
	message := req.Bytes()
	message = append(message, motorCodeBytes...)
	message = append(message, directionBytes...)
	message = append(message, speedBytes...)
	message = append(message, posBytes...)
	data, err := d.SendAck2(
		message,
		MotorMoveRelativeUnit.RecResp(),
		MotorMoveRelativeUnit.RecResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = binary.BigEndian.Uint16(data[3:5])
	return resp, nil
}

func (d *Dao) ReadOxygenConc() (resp interface{}, err error) {
	req := SensorOxygenConcUnit.Request()
	resp, err = d.Send(req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return resp, nil
}

func (d *Dao) Send(
	message []byte,
) ([]byte, error) {
	return d.UsbcanClient.Send(
		message,
		responseNil(),
		0,
		responseNil(),
		0,
	)
}

func (d *Dao) SendAck1(
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.UsbcanClient.Send(
		message,
		recResp,
		1,
		comResp,
		1,
	)
}

func (d *Dao) SendAck2(
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.UsbcanClient.Send(
		message,
		recResp,
		2,
		comResp,
		2,
	)
}

func stringToUint16Bytes(inputString string) (output []byte, isNegtive bool, err error) {
	input, err := strconv.Atoi(inputString)
	if err != nil {
		return output, isNegtive, err
	}
	if input < 0 {
		input = -input
		isNegtive = true
	}
	if input > math.MaxUint16 {
		return output, isNegtive, fmt.Errorf("%v overflows uint16", input)
	}
	var buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, uint16(input))
	if err != nil {
		return output, isNegtive, err
	}
	output = buf.Bytes()
	if len(output) != 2 {
		return output, isNegtive,
			fmt.Errorf("unexpected length %v of bytes %v", len(output), output)
	}
	return output, isNegtive, err
}

func stringToUint8Bytes(inputString string) (output []byte, err error) {
	input, err := strconv.Atoi(inputString)
	if err != nil {
		return output, err
	}
	if input > math.MaxUint8 {
		return output, fmt.Errorf("%v overflows uint8", input)
	}
	var buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, uint8(input))
	if err != nil {
		return output, err
	}
	output = buf.Bytes()
	if len(output) != 1 {
		return output,
			fmt.Errorf("unexpected length %v of bytes %v", len(output), output)
	}
	return output, err
}

func positionBytes(position string) ([]byte, []byte, error) {
	direction := []byte{0x00}
	posBytes, negtive, err := stringToUint16Bytes(position)
	if err != nil {
		return posBytes, direction, err
	}
	if negtive {
		direction = []byte{0x01}
	}
	return posBytes, direction, nil

}
