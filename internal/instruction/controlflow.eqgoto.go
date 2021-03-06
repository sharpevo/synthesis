package instruction

import (
	"fmt"
	"strconv"
	"synthesis/internal/dao"
	"synthesis/internal/interpreter/vrb"
)

func init() {
	dao.InstructionMap.Set("EQGOTO", InstructionControlFlowEqualGoto{})
}

type InstructionControlFlowEqualGoto struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowEqualGoto) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	index, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}
	cm, found := i.Env.Get("SYS_CMP")
	if !found {
		return resp, fmt.Errorf("failed to load variable CMP")
	}
	resp = fmt.Sprintf("condition 'equal to' check failed and continue")
	cm.Lock()
	v, _ := i.GetVarLockless(cm, "SYS_CMP")
	isEqual := v.GetValue() == vrb.EQUAL
	cm.Unlock()
	if isEqual {
		i.Goto(index)
		resp = fmt.Sprintf(
			"condition 'equal to' is satisfied and go to %d", index)
	}
	return resp, nil
}
