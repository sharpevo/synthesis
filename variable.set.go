package instruction

import (
	"fmt"
	"posam/interpreter/vrb"
)

type InstructionVariableSet struct {
	Instruction
}

func (i *InstructionVariableSet) Execute(args ...string) (resp interface{}, err error) {
	if len(args) <= 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name := args[0]
	value := args[1]
	variable, err := vrb.NewVariable(name, value)
	if err != nil {
		return resp, err
	}
	resp = i.Env.Set(name, variable)
	return fmt.Sprintf(
		"variable %q is set to %v",
		name,
		value,
	), nil
}
