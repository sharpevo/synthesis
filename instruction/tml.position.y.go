package instruction

import (
	"fmt"
	"synthesis/dao/aoztech"
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
	cm, err := i.ParseVariable(args[1]) // overwritten
	if err != nil {
		return resp, err
	}
	positionFloat := instance.TMLClient.PosY()
	positionBigFloat, err := i.PositionToBigFloat(positionFloat)
	if err != nil {
		return resp, err
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[1])
	variable.SetValue(positionBigFloat)
	cm.Unlock()
	resp = fmt.Sprintf("position y: %v", positionFloat)
	return resp, nil
}
