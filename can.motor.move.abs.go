package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("MOVEABS", InstructionCANMotorMoveAbsolute{})
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
	variable, err := i.ParseIntVariable(args[1])
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsolute(
		args[2],
		args[3],
	)
	if err != nil {
		return resp, err
	}
	variable.Value = resp
	return resp, err
}
