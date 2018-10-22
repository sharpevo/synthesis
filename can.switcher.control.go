package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("SWITCH", InstructionCANSwitcherControl{})
}

type InstructionCANSwitcherControl struct {
	InstructionCAN
}

func (i *InstructionCANSwitcherControl) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
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
	data, err := i.ParseInt(args[2])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ControlSwitcher(
		data,
	)
	if err != nil {
		return resp, err
	}
	variable.Value = resp
	return resp, err
}
