package instruction

import (
	"fmt"
	"synthesis/dao"
)

func init() {
	dao.InstructionMap.Set("RETURN", InstructionControlFlowReturn{})
}

type InstructionControlFlowReturn struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowReturn) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	i.Goto(-1)
	resp = fmt.Sprintf("return")
	return resp, nil
}
