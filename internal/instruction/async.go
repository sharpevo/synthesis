package instruction

import (
	"synthesis/internal/dao"
)

func init() {
	dao.InstructionMap.Set("ASYNC", InstructionAsync{})
}

type InstructionAsync struct {
	Instruction
}
