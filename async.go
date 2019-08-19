package instruction

import (
	"posam/dao"
)

func init() {
	dao.InstructionMap.Set("ASYNC", InstructionAsync{})
}

type InstructionAsync struct {
	Instruction
}
