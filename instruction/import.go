package instruction

import (
	"posam/dao"
)

func init() {
	dao.InstructionMap.Set("IMPORT", InstructionImport{})
}

type InstructionImport struct {
	Instruction
}
