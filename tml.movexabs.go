package instruction

import (
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEXABS", InstructionTMLMoveXABS{})
}

type InstructionTMLMoveXABS struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveXABS) Execute(args ...string) (resp interface{}, err error) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsByAxis(
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
