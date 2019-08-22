package aoztech

import (
	"fmt"
	"strconv"
	"sync"
	"synthesis/internal/dao"
	"synthesis/internal/protocol/tml"
	"synthesis/pkg/concurrentmap"
)

const (
	NAME        = "AOZTECH"
	DEVICE_NAME = "DEVICE_NAME"
	BAUD_RATE   = "BAUD_RATE"
	IDNAME      = DEVICE_NAME

	AXIS_X_ID         = "AXIS_X_ID"
	AXIS_X_SETUP_FILE = "AXIS_X_SETUP_FILE"
	AXIS_Y_ID         = "AXIS_Y_ID"
	AXIS_Y_SETUP_FILE = "AXIS_Y_SETUP_FILE"
)

var CONN_ATTRIBUTES = []string{
	DEVICE_NAME,
	BAUD_RATE,
	AXIS_X_ID,
	AXIS_X_SETUP_FILE,
	AXIS_Y_ID,
	AXIS_Y_SETUP_FILE,
}

var InstructionMap *dao.InstructionMapt

var deviceMap *concurrentmap.ConcurrentMap

func init() {
	InstructionMap = dao.NewInstructionMap()
	ResetInstance()
}

type Dao struct {
	id string
	sync.Mutex
	TMLClient *tml.Client
}

func NewDao(
	name string,
	baud string,
	axisXID string,
	axisXSetupFile string,
	axisYID string,
	axisYSetupFile string,
) (dao *Dao, err error) {
	baudInt, err := strconv.Atoi(baud)
	if err != nil {
		return dao, err
	}

	axisXIDInt, err := strconv.Atoi(axisXID)
	if err != nil {
		return dao, err
	}

	axisYIDInt, err := strconv.Atoi(axisYID)
	if err != nil {
		return dao, err
	}

	client, err := tml.NewClient(
		name,
		baudInt,
		axisXIDInt,
		axisXSetupFile,
		axisYIDInt,
		axisYSetupFile,
	)
	if err != nil {
		return dao, err
	}
	dao = &Dao{
		TMLClient: client,
	}

	if err = dao.SetID(name); err != nil {
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
	tml.ResetInstance()
}

func Instance(channel string) *Dao {
	if device, ok := deviceMap.Get(channel); ok {
		return device.(*Dao)
	}
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

func (d *Dao) Position() (float64, float64) {
	d.Lock()
	defer d.Unlock()
	return d.TMLClient.PosX(), d.TMLClient.PosY()
}

func (d *Dao) MoveRelByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) (resp interface{}, err error) {
	if err := d.TMLClient.MoveRelByAxis(
		axisID,
		pos,
		speed,
		accel,
	); err != nil {
		return resp, err
	}
	return fmt.Sprintf("moved by %v", pos), nil
}

func (d *Dao) MoveAbsByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) (resp interface{}, err error) {
	if err := d.TMLClient.MoveAbsByAxis(
		axisID,
		pos,
		speed,
		accel,
	); err != nil {
		return resp, err
	}
	return fmt.Sprintf("moved to %v", pos), nil
}

func (d *Dao) MoveRel(
	posx float64,
	posy float64,
	speed float64,
	accel float64,
) (resp interface{}, err error) {
	if err := d.TMLClient.MoveRel(
		posx,
		posy,
		speed,
		accel,
	); err != nil {
		return resp, err
	}
	return fmt.Sprintf("moved by (%v, %v)", posx, posy), nil
}

func (d *Dao) MoveAbs(
	posx float64,
	posy float64,
	speed float64,
	accel float64,
) (resp interface{}, err error) {
	if err := d.TMLClient.MoveAbs(
		posx,
		posy,
		speed,
		accel,
	); err != nil {
		return resp, err
	}
	return fmt.Sprintf("moved to (%v, %v)", posx, posy), nil
}
