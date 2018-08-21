package instruction

import (
	"fmt"
)

type InstructionDivisionFloat64 struct {
	InstructionArithmetic
}

func (i *InstructionDivisionFloat64) Execute(args ...string) (resp interface{}, err error) {
	variable, v1, v2, err := i.ParseObjects(args[0], args[1])
	if err != nil {
		return resp, err
	}
	v1.Quo(v1, v2)
	if v1.IsInf() {
		return resp, fmt.Errorf("inf quotition")
	}
	variable.Value = v1
	return v1.String(), nil
}
