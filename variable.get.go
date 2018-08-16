package instruction

import (
	"fmt"
)

type InstructionVariableGet struct {
	Instruction
}

func (i *InstructionVariableGet) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	variable, _ := i.ParseVariable(args[0])
	resp = variable.Value
	return resp, nil
}
