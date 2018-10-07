package instruction

import (
	"fmt"
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("POSITIONY", InstructionTMLPositionY{})
}

type InstructionTMLPositionY struct {
	InstructionTMLPosition
}

func (i *InstructionTMLPositionY) Execute(args ...string) (
	resp interface{},
	err error,
) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	instance, err := i.ParseDevice(args[0])
	if err != nil {
		return resp, err
	}
	variable, err := i.ParseVariable(args[1]) // overwritten
	if err != nil {
		return resp, err
	}
	positionFloat := instance.TMLClient.PosY
	positionBigFloat, err := i.PositionToBigFloat(positionFloat)
	if err != nil {
		return resp, err
	}
	variable.Value = positionBigFloat
	resp = fmt.Sprintf("position y: %v", variable.Value)
	return resp, nil
}
