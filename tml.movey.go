package instruction

import (
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEY", InstructionTMLMoveY{})
}

type InstructionTMLMoveY struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveY) Execute(args ...string) (resp interface{}, err error) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRelByAxis(
		instance.TMLClient.AxisYID,
		pos,
		speed,
		accel,
	)
	if err != nil {
		return resp, err
	}
	return
}
