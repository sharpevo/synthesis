package instruction

import (
	"fmt"
	"synthesis/internal/dao"
	"synthesis/internal/interpreter/vrb"
	"strconv"
)

func init() {
	dao.InstructionMap.Set("NEGOTO", InstructionControlFlowNotEqualGoto{})
}

type InstructionControlFlowNotEqualGoto struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowNotEqualGoto) Execute(args ...string) (resp interface{}, err error) {
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
	resp = fmt.Sprintf("condition 'not equal to' check failed and continue")
	cm.Lock()
	v, _ := i.GetVarLockless(cm, "SYS_CMP")
	isUnequal := v.GetValue() == vrb.UNEQUAL
	cm.Unlock()
	if isUnequal {
		i.Goto(index)
		resp = fmt.Sprintf(
			"condition 'not equal to' is satisfied and go to %d", index)
	}
	return resp, nil
}
