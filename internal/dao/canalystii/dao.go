package canalystii

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"synthesis/internal/dao"
	"synthesis/internal/protocol/usbcan"
	"synthesis/pkg/concurrentmap"
)

// Constants manage strings as the captions in the GUI. Note that the value of IDNAME is different according to the type of devices.
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

// CONN_ATTRIBUTES is a collection that manage attributes listed in the Deivce
// tree
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

// InstructionMap holds instructions registered in the instruction package.
var InstructionMap *dao.InstructionMapt
var deviceMap *concurrentmap.ConcurrentMap

func init() {
	InstructionMap = dao.NewInstructionMap()
	ResetInstance()
}

type Dao struct {
	_id          string
	usbcanClient *usbcan.Client
}

// NewDao is the constructor of canalystii.Dao, initializes Dao instance in
// singleton as well as UsbClient. Frame ID is exposed as DevID in the argument
// list.
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
		usbcanClient: usbcanClient,
	}
	if err = dao.setID(fmt.Sprintf("%v", devIDInt)); err != nil {
		return dao, err
	}
	return dao, nil
}

func addInstance(dao *Dao) {
	deviceMap.Set(dao.id(), dao)
}

func ResetInstance() {
	deviceMap = concurrentmap.NewConcurrentMap()
	usbcan.ResetInstance()
}

// Instance returns an instance of canalyst DAO by id, even id is an empty
// string.
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

func (d *Dao) id() string {
	return d._id
}

func (d *Dao) setID(id string) error {
	if _, ok := deviceMap.Get(id); ok {
		return fmt.Errorf("ID %q is duplicated", id)
	}
	d._id = id
	addInstance(d)
	return nil
}

// MoveRelative sends data of relative motion to the motor.
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
	output, err := sendAck2(d,
		composeBytes(
			req.Bytes(),
			motorCodeBytes,
			directionBytes,
			speedBytes,
			posBytes,
		),
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

// MoveAbsolute sends data of absolute motion to the motor.
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
	output, err := sendAck2(d,
		composeBytes(
			req.Bytes(),
			motorCodeBytes,
			posBytes,
			make([]byte, 3),
		),
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

// ResetMotor sends data to reset the motor
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
	output, err := sendAck2(d,
		composeBytes(
			req.Bytes(),
			motorCodeBytes,
			directionBytes,
			make([]byte, 4),
		),
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

// ControlSwitcher sends data to the switcher.
func (d *Dao) ControlSwitcher(
	data int,
) (resp interface{}, err error) {
	dataBytes, err := uint16Bytes(data)
	if err != nil {
		return resp, err
	}
	req := SwitcherControlUnit.Request()
	output, err := send(d,
		composeBytes(
			req.Bytes(),
			dataBytes,
			make([]byte, 4),
		))
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output
	return resp, nil
}

// ControlSwitcherAdvanced sends data to multiple switchers.
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
	output, err := sendAck6(d,
		composeBytes(
			req.Bytes(),
			dataBytes,
			speedBytes,
			countBytes,
			make([]byte, 1),
		),
		SwitcherControlAdvancedUnit.RecResp(),
		SwitcherControlAdvancedUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp = output
	return resp, nil
}

// ReadHumiture sends data to require humiture from sensors. The first value in
// the slice it returned is humidity, and the other one is temperature.
func (d *Dao) ReadHumiture() (resp interface{}, err error) {
	req := SensorHumitureUnit.Request()
	output, err := send(d, req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}
	humidity := binary.BigEndian.Uint16(output[1:3])
	temperature := binary.BigEndian.Uint16(output[3:6])
	resp = []float64{dividedByTen(humidity), dividedByTen(temperature)}
	return resp, nil
}

// ReadOxygenConc sends data to require concentration of oxygen from sensors.
func (d *Dao) ReadOxygenConc() (resp interface{}, err error) {
	req := SensorOxygenConcUnit.Request()
	output, err := send(d, req.Bytes())
	if err != nil {
		log.Println(err)
		return resp, err
	}
	conc := binary.BigEndian.Uint16(output[3:5])
	resp = dividedByTen(conc)
	return resp, nil
}

// ReadPressure sends data to require pressure from sensors. It returs the
// decimal value without convertion to the relative value based on the voltage.
func (d *Dao) ReadPressure(device int) (resp interface{}, err error) {
	deviceBytes, err := uint8Bytes(device)
	if err != nil {
		return resp, err
	}
	req := SensorPressureUnit.Request()
	output, err := send(d,
		composeBytes(
			req.Bytes(),
			deviceBytes,
			make([]byte, 5),
		))
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

// WriteSystemRom sends data to setup the system parameter.
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
	output, err := send1(d,
		composeBytes(
			req.Bytes(),
			addressBytes,
			valueBytes,
			make([]byte, 2),
		),
		SystemRomWriteUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	val := binary.BigEndian.Uint16(output[2:4])
	resp = int64(val)
	return resp, nil
}

// ReadSystemRom sends data to require system parameter.
func (d *Dao) ReadSystemRom(
	address int,
) (resp interface{}, err error) {
	addressBytes, err := uint16Bytes(address)
	if err != nil {
		return resp, err
	}
	req := SystemRomWriteUnit.Request()
	output, err := send1(d,
		composeBytes(
			req.Bytes(),
			addressBytes,
			make([]byte, 4),
		),
		SystemRomReadUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	val := binary.BigEndian.Uint16(output[2:4])
	resp = int64(val)
	return resp, nil
}

var send = func(
	d *Dao,
	message []byte,
) ([]byte, error) {
	return d.usbcanClient.Send(
		message,
		responseNil(),
		0,
		responseNil(),
		0,
	)
}

var send1 = func(
	d *Dao,
	message []byte,
	comResp []byte,
) ([]byte, error) {
	return d.usbcanClient.Send(
		message,
		responseNil(),
		0,
		comResp,
		1,
	)
}

var sendAck1 = func(
	d *Dao,
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.usbcanClient.Send(
		message,
		recResp,
		1,
		comResp,
		1,
	)
}

var sendAck2 = func(
	d *Dao,
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.usbcanClient.Send(
		message,
		recResp,
		2,
		comResp,
		2,
	)
}

var sendAck6 = func(
	d *Dao,
	message []byte,
	recResp []byte,
	comResp []byte,
) ([]byte, error) {
	return d.usbcanClient.Send(
		message,
		recResp,
		6,
		comResp,
		6,
	)
}

func uint16Bytes(input int) (output []byte, err error) {
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

func dividedByTen(input uint16) float64 {
	return float64(input) / 10.0
}

func composeBytes(target []byte, slices ...[]byte) []byte {
	for _, slice := range slices {
		target = append(target, slice...)
	}
	return target
}
