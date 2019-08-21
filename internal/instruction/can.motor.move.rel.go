package instruction

import (
	"fmt"
	"synthesis/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("CANMOVEREL", InstructionCANMotorMoveRelative{})
}

type InstructionCANMotorMoveRelative struct {
	InstructionCAN
}

func (i *InstructionCANMotorMoveRelative) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 6 {
		return resp, fmt.Errorf("not enough arguments")
	}
	instance, err := i.ParseDevice(args[0])
	if err != nil {
		return resp, err
	}
	cm, err := i.ParseIntVariable(args[1])
	if err != nil {
		return resp, err
	}
	motorCode, err := i.ParseInt(args[2])
	if err != nil {
		return resp, err
	}
	direction, err := i.ParseInt(args[3])
	if err != nil {
		return resp, err
	}
	speed, err := i.ParseInt(args[4])
	if err != nil {
		return resp, err
	}
	position, err := i.ParseInt(args[5])
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRelative(
		motorCode,
		direction,
		speed,
		position,
	)
	if err != nil {
		return resp, err
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[1])
	variable.SetValue(resp)
	cm.Unlock()
	return resp, err
}
