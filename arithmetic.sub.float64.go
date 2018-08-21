package instruction

import ()

type InstructionSubtractionFloat64 struct {
	InstructionArithmetic
}

func (i *InstructionSubtractionFloat64) Execute(args ...string) (resp interface{}, err error) {
	variable, v1, v2, err := i.ParseObjects(args[0], args[1])
	if err != nil {
		return resp, err
	}
	v1.Sub(v1, v2)
	variable.Value = v1
	return v1.String(), nil
}
