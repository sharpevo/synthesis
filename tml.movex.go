package instruction

import (
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEX", InstructionTMLMoveX{})
}

type InstructionTMLMoveX struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveX) Execute(args ...string) (resp interface{}, err error) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRelByAxis(
		instance.TMLClient.AxisXID,
		pos,
		speed,
		accel,
	)
	if err != nil {
		return resp, err
	}
	return
}
