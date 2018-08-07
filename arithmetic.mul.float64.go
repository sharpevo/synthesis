package instruction

import ()

type InstructionMultiplicationFloat64 struct {
	InstructionArithmetic
}

func (i *InstructionMultiplicationFloat64) Execute(args ...string) (resp interface{}, err error) {
	variable, v1, v2, err := i.ParseObjects(args[0], args[1])
	if err != nil {
		return resp, err
	}
	v1.Mul(v1, v2)
	variable.Value = v1.String()
	return v1.String(), nil
}
