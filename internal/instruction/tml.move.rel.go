package instruction

import (
	"synthesis/internal/dao/aoztech"
)

func init() {
	aoztech.InstructionMap.Set("MOVEREL", InstructionTMLMoveRel{})
}

type InstructionTMLMoveRel struct {
	InstructionTMLMove
}

func (i *InstructionTMLMoveRel) Execute(args ...string) (
	resp interface{},
	err error,
) {
	instance, posx, posy, speed, accel, err := i.InitializeDual(args...)
	if err != nil {
		return resp, err
	}
	resp, err = instance.MoveRel(
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
