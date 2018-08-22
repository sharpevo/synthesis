package instruction

import (
	"fmt"
	"posam/interpreter/vrb"
	"strconv"
)

type InstructionControlFlowLessThanGoto struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowLessThanGoto) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	index, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}
	v, found := i.Env.Get("SYS_CMP")
	if !found {
		return resp, fmt.Errorf("failed to load variable CMP")
	}
	resp = fmt.Sprintf("condition cheking passed and continue")
	if v.Value == vrb.LESS {
		i.Goto(index)
		resp = fmt.Sprintf("condition satisfied and goto %d", index)
	}
	return resp, nil
}
