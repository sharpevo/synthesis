package instruction

import (
	"synthesis/internal/dao"
)

func init() {
	dao.InstructionMap.Set("IMPORT", InstructionImport{})
}

type InstructionImport struct {
	Instruction
}
