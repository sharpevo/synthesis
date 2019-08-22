package instruction

import (
	"fmt"
	"synthesis/internal/dao/aoztech"
	"strconv"
)

type InstructionTML struct {
	Instruction
}

func (i *InstructionTML) ParseFloat(input string) (output float64, err error) {
	cm, found := i.Env.Get(input)
	if !found {
		output, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return output, err
		}
	} else {
		cm.Lock()
		outputVar, _ := i.GetVarLockless(cm, input)
		output, err = strconv.ParseFloat(fmt.Sprintf("%v", outputVar.GetValue()), 64)
		cm.Unlock()
		if err != nil {
			return output,
				fmt.Errorf(
					"failed to parse variable %q to float: %s",
					input,
					err.Error(),
				)
		}
	}
	return output, nil
}

func (i *InstructionTML) ParseDevice(input string) (
	instance *aoztech.Dao,
	err error,
) {
	cm, found := i.Env.Get(input)
	if !found {
		return instance,
			fmt.Errorf("device %q is not defined", input)
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, input)
	deviceName := variable.GetValue().(string)
	cm.Unlock()
	instance = aoztech.Instance(deviceName)
	if instance == nil {
		return instance,
			fmt.Errorf("device %q is not initialized", input)
	}
	return instance, nil
}
