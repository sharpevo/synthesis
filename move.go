package instruction

import (
	"fmt"
)

type InstructionMove struct {
	Instruction
}

func (i *InstructionMove) Execute(args ...string) (interface{}, error) {
	if i.isMovable() {
		result := fmt.Sprintf("Movable: %s", args[0])
		return result, nil
	} else {
		fmt.Println("Can not move")
		return i.Instruction.Execute(args...)
	}
}

func (i *InstructionMove) isMovable() bool {
	return true
}
