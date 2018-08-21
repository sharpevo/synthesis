package instruction

import (
	"fmt"
	"math/big"
	"posam/interpreter/vrb"
)

type InstructionArithmetic struct {
	Instruction
}

func (i *InstructionArithmetic) ParseObjects(arg1 string, arg2 string) (
	variable *vrb.Variable,
	v1 *big.Float,
	v2 *big.Float,
	err error,
) {
	variable, found := i.Env.Get(arg1)
	if !found {
		return variable, v1, v2, fmt.Errorf("Invalid variable %q", arg1)
	}

	v1, err = i.GetBigFloat64(arg1)
	if err != nil {
		return variable, v1, v2, err
	}

	v2, err = i.GetBigFloat64(arg2)
	if err != nil {
		return variable, v1, v2, err
	}

	return variable, v1, v2, nil
}

func (i *InstructionArithmetic) GetBigFloat64(
	input interface{},
) (*big.Float, error) {
	switch v := input.(type) {
	case string:
		variable, found := i.Env.Get(v)
		if found {
			return i.GetBigFloat64(variable)
		} else {
			floatv, _, err := big.ParseFloat(v, 10, 53, big.ToNearestEven)
			if err != nil {
				return big.NewFloat(0), err
			}
			return floatv, nil
		}
	case *vrb.Variable:
		return i.GetBigFloat64(v.Value)
	case *big.Float:
		return v, nil
	default:
		return big.NewFloat(0), fmt.Errorf("invalid float %v", v)
	}
}
