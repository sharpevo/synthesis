package instruction

import (
	"fmt"
	"synthesis/internal/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("CANMOVEABS", InstructionCANMotorMoveAbsolute{})
}

type InstructionCANMotorMoveAbsolute struct {
	InstructionCAN
}

func (i *InstructionCANMotorMoveAbsolute) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 4 {
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
	position, err := i.ParseInt(args[3])
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsolute(
		motorCode,
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
