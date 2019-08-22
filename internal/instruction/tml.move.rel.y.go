package instruction

import (
	"synthesis/internal/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVERELY", InstructionTMLMoveRelY{})
}

type InstructionTMLMoveRelY struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveRelY) Execute(args ...string) (
	resp interface{},
	err error,
) {
	instance, pos, speed, accel, err := i.Initialize(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRelByAxis(
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
