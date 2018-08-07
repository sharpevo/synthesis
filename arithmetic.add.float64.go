package instruction

import (
	//"math/big"
	"fmt"
	"posam/interpreter"
)

type InstructionAddFloat64 struct {
	InstructionArithmetic
}

func (i *InstructionAddFloat64) Execute(args ...string) (resp interface{}, err error) {

	v, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("Invalid variable %q", args[0])
	}
	variable, ok := v.(*interpreter.Variable)
	if !ok {
		return resp, fmt.Errorf("Invalid type of variable %q", args[0])
	}

	v1, err := i.GetBigFloat64(args[0])
	if err != nil {
		return resp, err
	}

	v2, err := i.GetBigFloat64(args[1])
	if err != nil {
		return resp, err
	}
	v1.Add(v1, v2)

	variable.Value = v1.String()
	return v1.String(), nil
}
