package instruction

import (
	"posam/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVERELX", InstructionTMLMoveRelX{})
}

type InstructionTMLMoveRelX struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveRelX) Execute(args ...string) (
	resp interface{},
	err error,
) {
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
