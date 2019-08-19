package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("ROMREAD", InstructionCANSystemRomRead{})
}

type InstructionCANSystemRomRead struct {
	InstructionCAN
}

func (i *InstructionCANSystemRomRead) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
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
	address, err := i.ParseInt(args[2])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ReadSystemRom(
		address,
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
