package instruction

import (
	"synthesis/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEABSY", InstructionTMLMoveAbsY{})
}

type InstructionTMLMoveAbsY struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveAbsY) Execute(args ...string) (
	resp interface{},
	err error,
) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsByAxis(
		instance.TMLClient.AxisYID(),
		pos,
		speed,
		accel,
	)
	if err != nil {
		return resp, err
	}
	return
}
