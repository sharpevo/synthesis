package instruction

import (
	"fmt"
	"posam/dao/aoztech"
	"strconv"
)

type InstructionTML struct {
	Instruction
}

func (i *InstructionTML) ParseFloat(input string) (output float64, err error) {
	outputVar, found := i.Env.Get(input)
	if !found {
		output, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return output, err
		}
	} else {
		output, err = strconv.ParseFloat(fmt.Sprintf("%v", outputVar.Value), 64)
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
	variable, found := i.Env.Get(input)
	if !found {
		return instance,
			fmt.Errorf("device %q is not defined", input)
	}
	deviceName := variable.Value.(string)
	instance = aoztech.Instance(deviceName)
	if instance == nil {
		return instance,
			fmt.Errorf("device %q is not initialized", input)
	}
	return instance, nil
}
