package instruction

import (
	"fmt"
	"math/big"
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
	tempCM, err := i.ParseFloat64Variable(args[1])
	if err != nil {
		return resp, err
	}
	humiCM, err := i.ParseFloat64Variable(args[2])
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
	tempCM.Lock()
	tempVariable, _ := i.GetVarLockless(tempCM, args[1])
	tempVariable.SetValue(big.NewFloat(humiture[0]))
	tempCM.Unlock()
	humiCM.Lock()
	humiVariable, _ := i.GetVarLockless(humiCM, args[2])
	humiVariable.SetValue(big.NewFloat(humiture[1]))
	humiCM.Unlock()
	return resp, err
}
