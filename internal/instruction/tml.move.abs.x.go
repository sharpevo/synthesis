package instruction

import (
	"synthesis/internal/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEABSX", InstructionTMLMoveAbsX{})
}

type InstructionTMLMoveAbsX struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveAbsX) Execute(args ...string) (
	resp interface{},
	err error,
) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbsByAxis(
		instance.TMLClient.AxisXID(),
		pos,
		speed,
		accel,
	)
	if err != nil {
		return resp, err
	}
	return
}
