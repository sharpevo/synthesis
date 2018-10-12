package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("MOVEREL", InstructionCANMotorMoveRelative{})
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
	variable, err := i.ParseIntVariable(args[1])
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRelative(
		args[2],
		args[3],
		args[4],
		args[5],
	)
	if err != nil {
		return resp, err
	}
	variable.Value = resp
	return resp, err
}
