package instruction

import (
	"fmt"
	"math/big"
	"posam/interpreter"
)

type InstructionArithmetic struct {
	Instruction
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
	case *interpreter.Variable:
		return i.GetBigFloat64(v.Value)
	case *big.Float:
		return v, nil
	default:
		return big.NewFloat(0), fmt.Errorf("invalid float %v", v)
	}
}
