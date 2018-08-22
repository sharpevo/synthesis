package instruction

import (
	"fmt"
)

type InstructionControlFlowGoto struct {
	Instruction
}

func (i *InstructionControlFlowGoto) Goto(index int64) error {
	variable, found := i.Env.Get("SYS_CUR")
	if !found {
		return fmt.Errorf("invalid variable CUR")
	}
	variable.Value = index
	return nil
}
