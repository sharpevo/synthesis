package instruction

import (
	"synthesis/dao"
)

func init() {
	dao.InstructionMap.Set("ASYNC", InstructionAsync{})
}

type InstructionAsync struct {
	Instruction
}
