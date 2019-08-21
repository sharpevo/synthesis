package instruction

import (
	"fmt"
	"synthesis/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("OXYGENCONC", InstructionCANSensorOxygenConc{})
}

type InstructionCANSensorOxygenConc struct {
	InstructionCAN
}

func (i *InstructionCANSensorOxygenConc) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	instance, err := i.ParseDevice(args[0])
	if err != nil {
		return resp, err
	}
	cm, err := i.ParseFloat64Variable(args[1])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ReadOxygenConc()
	if err != nil {
		return resp, err
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[1])
	variable.SetValue(resp)
	cm.Unlock()
	return resp, err
}
