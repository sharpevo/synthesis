package instruction

import (
	"fmt"
	"synthesis/internal/dao"
	"strconv"
)

func init() {
	dao.InstructionMap.Set("LOOP", InstructionControlFlowLoop{})
}

type InstructionControlFlowLoop struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowLoop) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	cm, found := i.Env.Get(args[1])
	if !found {
		return resp, fmt.Errorf("variable %q is not defined", args[1])
	}

	cm.Lock()
	varTotal, _ := i.GetVarLockless(cm, args[1])
	total64 := varTotal.GetValue().(int64) - 1
	cm.Unlock() // Goto requires lock

	if total64 < 1 {
		return fmt.Sprintf("loop done"), nil
	}

	line, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}

	i.Goto(line)
	cm.Lock()
	varTotal.SetValue(total64)
	cm.Unlock()
	resp = fmt.Sprintf("%d left loops from line %d", total64, line)
	return resp, nil
}
