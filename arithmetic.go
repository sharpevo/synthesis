package instruction

import (
	"fmt"
	"math/big"
	"posam/interpreter/vrb"
	"strconv"
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

func (i *InstructionArithmetic) GetInt64(
	input interface{},
) (int64, error) {
	switch v := input.(type) {
	case string:
		variable, found := i.Env.Get(v)
		if found {
			return i.GetInt64(variable)
		} else {
			output, err := strconv.ParseInt(v, 0, 64)
			if err != nil {
				return 0, err
			}
			return output, nil
		}
	case *vrb.Variable:
		return i.GetInt64(v.Value)
	case int64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("invalid float %v", v)
	}
}

func (i *InstructionArithmetic) ParseOperands(arg1 string, arg2 string) (
	variable *vrb.Variable,
	v1 interface{},
	v2 interface{},
	err error,
) {
	variable, found := i.Env.Get(arg1)
	if !found {
		return variable, v1, v2, fmt.Errorf("Invalid variable %q", arg1)
	}
	switch variable.Type {
	case vrb.FLOAT:
		v1, err = i.GetBigFloat64(arg1)
		if err != nil {
			return variable, v1, v2, err
		}
		v2, err = i.GetBigFloat64(arg2)
		if err != nil {
			return variable, v1, v2, err
		}
	case vrb.INT:
		v1, err = i.GetInt64(arg1)
		if err != nil {
			return variable, v1, v2, err
		}
		v2, err = i.GetInt64(arg2)
		if err != nil {
			return variable, v1, v2, err
		}
	default:
		return variable, v1, v2,
			fmt.Errorf("invalid variable type %q", variable.Type)
	}
	return variable, v1, v2, nil
}
