package instruction

import (
	"fmt"
	"synthesis/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("SWITCHCOND", InstructionCANSwitcherControlAdvanced{})
}

type InstructionCANSwitcherControlAdvanced struct {
	InstructionCAN
}

func (i *InstructionCANSwitcherControlAdvanced) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 5 {
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
	data, err := i.ParseInt(args[2])
	if err != nil {
		return resp, err
	}
	speed, err := i.ParseInt(args[3])
	if err != nil {
		return resp, err
	}
	count, err := i.ParseInt(args[4])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ControlSwitcherAdvanced(
		data,
		speed,
		count,
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
