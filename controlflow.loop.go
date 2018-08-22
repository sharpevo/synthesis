package instruction

import (
	"fmt"
	"strconv"
)

type InstructionControlFlowLoop struct {
	InstructionControlFlowGoto
	count int
}

func (i *InstructionControlFlowLoop) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}

	line, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[0], err.Error())
	}

	total64, err := strconv.ParseInt(args[1], 0, 64)
	if err != nil {
		return resp, fmt.Errorf("invalid argument %q: %s", args[1], err.Error())
	}
	total := int(total64)
	if i.Count() == total {
		return fmt.Sprintf("loop done"), nil
	}
	i.Goto(line)
	i.SetCount(i.Count() + 1)

	resp = fmt.Sprintf("loop %d/%d from line %d", i.Count(), total, line)
	return resp, nil
}

func (i *InstructionControlFlowLoop) Count() int {
	return i.count
}

func (i *InstructionControlFlowLoop) SetCount(count int) {
	i.count = count
}
