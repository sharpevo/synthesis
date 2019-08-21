package instruction

import (
	"synthesis/dao"
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
