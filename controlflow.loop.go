package instruction

import (
	"fmt"
	"strconv"
)

type InstructionControlFlowLoop struct {
	InstructionControlFlowGoto
}

func (i *InstructionControlFlowLoop) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	varTotal, found := i.Env.Get(args[1])
	if !found {
		return resp, fmt.Errorf("variable %q is not defined", args[1])
	}

	total64 := varTotal.Value.(int64)

	if total64 == 0 {
		return fmt.Sprintf("loop done"), nil
	}

	line, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}

	i.Goto(line)
	varTotal.Value = total64 - 1
	resp = fmt.Sprintf("%d left loops from line %d", total64, line)
	return resp, nil
}
