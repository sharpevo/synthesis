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
	name := args[0]
	variable, found := i.Env.Get(name)
	if !found {
		resp = fmt.Sprintf("%s is not defined", name)
		return
	}
	resp = fmt.Sprintf("%v %s = %v", variable.Type, name, variable.Value)
	return resp, nil
}
