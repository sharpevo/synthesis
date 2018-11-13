package instruction

import (
	"fmt"
	"posam/dao/canalystii"
)

func init() {
	canalystii.InstructionMap.Set("HUMITURE", InstructionCANSensorHumiture{})
}

type InstructionCANSensorHumiture struct {
	InstructionCAN
}

func (i *InstructionCANSensorHumiture) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
		return resp, fmt.Errorf("not enough arguments")
	}
	instance, err := i.ParseDevice(args[0])
	if err != nil {
		return resp, err
	}
	tempVariable, err := i.ParseFloat64Variable(args[1])
	if err != nil {
		return resp, err
	}
	humiVariable, err := i.ParseFloat64Variable(args[2])
	if err != nil {
		return resp, err
	}
	resp, err = instance.ReadHumiture()
	if err != nil {
		return resp, err
	}
	humiture, ok := resp.([]float64)
	if !ok {
		return resp, fmt.Errorf("invalid humiture response %#v", resp)
	}
	tempVariable.Value = humiture[0]
	humiVariable.Value = humiture[1]
	return resp, err
}
