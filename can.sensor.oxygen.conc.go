package instruction

import (
	"fmt"
	"posam/dao/canalystii"
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
	variable, err := i.ParseFloat64Variable(args[1])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ReadOxygenConc()
	if err != nil {
		return resp, err
	}
	variable.Value = resp
	return resp, err
}
