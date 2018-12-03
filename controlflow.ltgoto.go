package instruction

import (
	"fmt"
	"posam/dao"
	"posam/interpreter/vrb"
	"strconv"
)

func init() {
	dao.InstructionMap.Set("LTGOTO", InstructionControlFlowLessThanGoto{})
}

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
	resp = fmt.Sprintf("condition 'less than' check failed and continue")
	if v.Value == vrb.LESS {
		i.Goto(index)
		resp = fmt.Sprintf(
			"condition 'less than' is satisfied and go to %d", index)
	}
	return resp, nil
}
