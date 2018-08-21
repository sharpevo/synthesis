package instruction

import (
	"fmt"
	"posam/interpreter/vrb"
)

type InstructionVariableCompare struct {
	Instruction
}

func (i *InstructionVariableCompare) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name1 := args[0]
	name2 := args[1]
	var1, err := i.ParseVariable(name1)
	if err != nil {
		return resp, err
	}
	var2, err := i.ParseVariable(name2)
	if err != nil {
		return resp, err
	}
	result, err := vrb.Compare(var1, var2)
	if err != nil {
		return resp, err
	}
	v, found := i.Env.Get("SYS_CMP")
	if !found {
		return resp, fmt.Errorf("failed to load variable CMP")
	}
	v.Value = result
	return result, nil
}
