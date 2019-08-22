package instruction

import (
	"fmt"
	"synthesis/internal/dao"
	"strconv"
)

func init() {
	dao.InstructionMap.Set("ERRGOTO", InstructionControlFlowErrGoto{})
}

type InstructionControlFlowErrGoto struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowErrGoto) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	index, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}
	cm, found := i.Env.Get("SYS_ERR")
	if !found {
		return resp, fmt.Errorf("failed to load variable ERR")
	}
	resp = fmt.Sprintf("error check failed and continue")
	cm.Lock()
	v, _ := i.GetVarLockless(cm, "SYS_ERR")
	isValid := v.GetValue().(string) != ""
	cm.Unlock()
	if isValid {
		i.Goto(index)
		resp = fmt.Sprintf("error detected and go to %d", index)
	}
	return resp, nil
}
