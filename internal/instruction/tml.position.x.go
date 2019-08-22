package instruction

import (
	"fmt"
	"synthesis/internal/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("POSITIONX", InstructionTMLPositionX{})
}

type InstructionTMLPositionX struct {
	InstructionTMLPosition
}

func (i *InstructionTMLPositionX) Execute(args ...string) (
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
	cm, err := i.ParseVariable(args[1]) // overwritten
	if err != nil {
		return resp, err
	}
	positionFloat := instance.TMLClient.PosX()
	positionBigFloat, err := i.PositionToBigFloat(positionFloat)
	if err != nil {
		return resp, err
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[1])
	variable.SetValue(positionBigFloat)
	cm.Unlock()
	resp = fmt.Sprintf("position x: %v", positionFloat)
	return resp, nil
}
