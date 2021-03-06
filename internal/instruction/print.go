package instruction

import (
	"synthesis/internal/dao"
	"strings"
)

func init() {
	dao.InstructionMap.Set("PRINT", InstructionPrint{})
}

type InstructionPrint struct {
	Instruction
}

func (c *InstructionPrint) Execute(args ...string) (interface{}, error) {
	message := strings.Join(args, " ")
	return message, nil
}
