package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

type InstructionCAN struct {
	Instruction
}

func (i *InstructionCAN) ParseDevice(input string) (
	instance *canalystii.Dao,
	err error,
) {
	variable, found := i.Env.Get(input)
	if !found {
		return instance,
			fmt.Errorf("device %q is not defined", input)
	}
	deviceID := fmt.Sprintf("%v", variable.Value)
	instance, err = canalystii.Instance(deviceID)
	if err != nil {
		return instance, err
	}
	return instance, nil
}
