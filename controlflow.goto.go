package instruction

import (
	"fmt"
	"strconv"
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

func (i *InstructionControlFlowGoto) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	line, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}
	i.Goto(line)
	resp = fmt.Sprintf("go to line %d", line)
	return resp, nil
}
