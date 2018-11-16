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
	motorCode int,
	direction int,
	speed int,
	position int,
) (resp interface{}, err error) {
	motorCodeBytes, err := uint8Bytes(motorCode)
	if err != nil {
		return resp, err
	}
	directionBytes, err := uint8Bytes(direction)
	if err != nil {
		return resp, err
	}
	speedBytes, err := uint16Bytes(speed)
	if err != nil {
		return resp, err
	}
	posBytes, err := uint16Bytes(position)
	if err != nil {
		return resp, err
	}
	req := MotorMoveRelativeUnit.Request()
	message := req.Bytes()
	message = append(message, motorCodeBytes...)
	message = append(message, directionBytes...)
	message = append(message, speedBytes...)
	message = append(message, posBytes...)
	output, err := d.SendAck2(
		message,
		MotorMoveRelativeUnit.RecResp(),
		MotorMoveRelativeUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = binary.BigEndian.Uint16(output[3:5])
	return resp, nil
}

func (d *Dao) MoveAbsolute(
	motorCode int,
	position int,
) (resp interface{}, err error) {
	motorCodeBytes, err := uint8Bytes(motorCode)
	if err != nil {
		return resp, err
	}
	posBytes, err := uint16Bytes(position)
	if err != nil {
		return resp, err
	}
	req := MotorMoveAbsoluteUnit.Request()
	message := req.Bytes()
	message = append(message, motorCodeBytes...)
	message = append(message, posBytes...)
	message = append(message, []byte{0x00, 0x00, 0x00}...)
	output, err := d.SendAck2(
		message,
		MotorMoveAbsoluteUnit.RecResp(),
		MotorMoveAbsoluteUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = binary.BigEndian.Uint16(output[3:5])
	return resp, nil
}

func (d *Dao) ResetMotor(
	motorCode int,
	direction int,
) (resp interface{}, err error) {
	motorCodeBytes, err := uint8Bytes(motorCode)
	if err != nil {
		return resp, err
	}
	directionBytes, err := uint8Bytes(direction)
	if err != nil {
		return resp, err
	}
	req := MotorResetUnit.Request()
	message := req.Bytes()
	message = append(message, motorCodeBytes...)
	message = append(message, directionBytes...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00}...)
	output, err := d.SendAck2(
		message,
		MotorResetUnit.RecResp(),
		MotorResetUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output
	return resp, nil
}

func (d *Dao) ControlSwitcher(
	data int,
) (resp interface{}, err error) {
	dataBytes, err := uint16Bytes(data)
	if err != nil {
		return resp, err
	}
	req := SwitcherControlUnit.Request()
	message := req.Bytes()
	message = append(message, dataBytes...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00}...)
	output, err := d.Send(
		message,
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output
	return resp, nil
}

func (d *Dao) ControlSwitcherAdvanced(
	data int,
	speed int,
	count int,
) (resp interface{}, err error) {
	dataBytes, err := uint16Bytes(data)
	if err != nil {
		return resp, err
	}
	speedBytes, err := uint8Bytes(speed)
	if err != nil {
		return resp, err
	}
	countBytes, err := uint16Bytes(count)
	if err != nil {
		return resp, err
	}
	req := SwitcherControlAdvancedUnit.Request()
	message := req.Bytes()
	message = append(message, dataBytes...)
	message = append(message, speedBytes...)
	message = append(message, countBytes...)
	message = append(message, 0x00)
	output, err := d.SendAck6(
		message,
		SwitcherControlUnit.RecResp(),
		SwitcherControlUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output
	return resp, nil
}

func (d *Dao) ReadHumiture() (resp interface{}, err error) {
	req := SensorHumitureUnit.Request()
	output, err := d.Send(req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}
	temperature := binary.BigEndian.Uint16(output[1:3])
	humidity := binary.BigEndian.Uint16(output[3:6])
	resp = []float64{divideTen(temperature), divideTen(humidity)}
	return resp, nil
}

func (d *Dao) ReadOxygenConc() (resp interface{}, err error) {
	req := SensorOxygenConcUnit.Request()
	output, err := d.Send(req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}
	conc := binary.BigEndian.Uint16(output[3:5])
	//conc := binary.BigEndian.Uint16(output[2:4]) // canalystii
	//conc := binary.BigEndian.Uint16(output[1:3]) // usbcan-2c
	resp = divideTen(conc)
	return resp, nil
}

func (d *Dao) ReadPressure(device int) (resp interface{}, err error) {
	deviceBytes, err := uint8Bytes(device)
	if err != nil {
		return resp, err
	}
	req := SensorPressureUnit.Request()
	message := req.Bytes()
	message = append(message, deviceBytes...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00}...)
	output, err := d.Send(message)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	if output[2] == 0xff {
		return resp, fmt.Errorf("invalid pressure device '%v'", device)
	}
	voltageDec := binary.BigEndian.Uint16(output[2:4])
	resp = int64(voltageDec)
	return resp, nil
}

func (d *Dao) WriteSystemRom(
	address int,
	value int,
) (resp interface{}, err error) {
	addressBytes, err := uint16Bytes(address)
	if err != nil {
		return resp, err
	}
	valueBytes, err := uint16Bytes(value)
	if err != nil {
		return resp, err
	}
	req := SystemRomWriteUnit.Request()
	message := req.Bytes()
	message = append(message, addressBytes...)
	message = append(message, valueBytes...)
	message = append(message, []byte{0x00, 0x00}...)
	output, err := d.Send1(
		message,
		SystemRomWriteUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output[2:4]
	return resp, nil
}

func (d *Dao) ReadSystemRom(
	address int,
) (resp interface{}, err error) {
	addressBytes, err := uint16Bytes(address)
	if err != nil {
		return resp, err
	}
	req := SystemRomWriteUnit.Request()
	message := req.Bytes()
	message = append(message, addressBytes...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00}...)
	output, err := d.Send1(
		message,
		SystemRomReadUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output[2:4]
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

func (d *Dao) Send1(
	message []byte,
	comResp []byte,
) ([]byte, error) {
	return d.UsbcanClient.Send(
		message,
		responseNil(),
		0,
		comResp,
		1,
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

func (d *Dao) SendAck6(
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.UsbcanClient.Send(
		message,
		recResp,
		6,
		comResp,
		6,
	)
}

func uint16Bytes(input int) (output []byte, err error) {
	//if input < 0 {
	//input = -input
	//isNegtive = true
	//}
	if input > math.MaxUint16 {
		return output, fmt.Errorf("%v overflows uint16", input)
	}
	var buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, uint16(input))
	if err != nil {
		return output, err
	}
	output = buf.Bytes()
	if len(output) != 2 {
		return output,
			fmt.Errorf("unexpected length %v of bytes %v", len(output), output)
	}
	return output, err
}

func uint8Bytes(input int) (output []byte, err error) {
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

func divideTen(input uint16) float64 {
	return float64(input) / 10.0
}
