package instruction

import (
	"synthesis/internal/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEABS", InstructionTMLMoveAbs{})
}

type InstructionTMLMoveAbs struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveAbs) Execute(args ...string) (
	resp interface{},
	err error,
) {
	instance, posx, posy, speed, accel, err := i.InitializeDual(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveAbs(
		posx,
		posy,
		speed,
		accel,
	)
	if err != nil {
		return resp, err
	}
	return
}
