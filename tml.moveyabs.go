package instruction

import (
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEYABS", InstructionTMLMoveYABS{})
}

type InstructionTMLMoveYABS struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveYABS) Execute(args ...string) (resp interface{}, err error) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsByAxis(
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
