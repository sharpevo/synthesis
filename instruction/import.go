package instruction

import (
	"synthesis/dao"
)

func init() {
	dao.InstructionMap.Set("IMPORT", InstructionImport{})
}

type InstructionImport struct {
	Instruction
}
