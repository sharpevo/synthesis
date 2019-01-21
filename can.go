package instruction

import (
	"fmt"
	"posam/dao/canalystii"
	"posam/interpreter/vrb"
)

type InstructionCAN struct {
	Instruction
}

func (i *InstructionCAN) ParseDevice(input string) (
	instance *canalystii.Dao,
	err error,
) {
	cm, found := i.Env.Get(input)
	if !found {
		return instance,
			fmt.Errorf("device %q is not defined", input)
	}
	cm.Lock()
	variablei, _ := cm.GetLockless(input)
	variable, _ := variablei.(*vrb.Variable)
	deviceID := fmt.Sprintf("%v", variable.GetValue())
	cm.Unlock()
	instance, err = canalystii.Instance(deviceID)
	if err != nil {
		return instance, err
	}
	return instance, nil
}
