package instruction

import (
	"fmt"
	"math/big"
	"strconv"
	"synthesis/internal/interpreter/vrb"
	"synthesis/pkg/concurrentmap"
)

type InstructionArithmetic struct {
	Instruction
}

func (i *InstructionArithmetic) ParseObjects(arg1 string, arg2 string) (
	cm *concurrentmap.ConcurrentMap,
	v1 *big.Float,
	v2 *big.Float,
	err error,
) {
	cm, found := i.Env.Get(arg1)
	if !found {
		return cm, v1, v2, fmt.Errorf("Invalid variable %q", arg1)
	}

	v1, err = i.GetBigFloat64(arg1)
	if err != nil {
		return cm, v1, v2, err
	}

	v2, err = i.GetBigFloat64(arg2)
	if err != nil {
		return cm, v1, v2, err
	}

	return cm, v1, v2, nil
}

func (i *InstructionArithmetic) GetBigFloat64(
	input interface{},
) (*big.Float, error) {
	switch v := input.(type) {
	case string:
		cm, found := i.Env.Get(v)
		if found {
			cm.Lock()
			defer cm.Unlock()
			variablei, _ := cm.GetLockless(v)
			variable, _ := variablei.(*vrb.Variable)
			v, ok := variable.GetValue().(*big.Float)
			if !ok {
				return big.NewFloat(0), fmt.Errorf("invalid float %v(%T)", v, v)
			}
			return v, nil
		} else {
			floatv, _, err := big.ParseFloat(v, 10, 53, big.ToNearestEven)
			if err != nil {
				return big.NewFloat(0), err
			}
			return floatv, nil
		}
	case *big.Float:
		return v, nil
	default:
		return big.NewFloat(0), fmt.Errorf("invalid float %v(%T)", v, v)
	}
}

func (i *InstructionArithmetic) GetInt64(
	input interface{},
) (int64, error) {
	switch v := input.(type) {
	case string:
		cm, found := i.Env.Get(v)
		if found {
			cm.Lock()
			defer cm.Unlock()
			variablei, _ := cm.GetLockless(v)
			variable, _ := variablei.(*vrb.Variable)
			vi := variable.GetValue()
			v, ok := vi.(int64)
			if !ok {
				return 0, fmt.Errorf("invalid int %v (%T)", vi, vi)
			}
			return v, nil
		} else {
			output, err := strconv.ParseInt(v, 0, 64)
			if err != nil {
				return 0, err
			}
			return output, nil
		}
	case int64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("invalid int %v", v)
	}
}

func (i *InstructionArithmetic) ParseOperands(arg1 string, arg2 string) (
	cm *concurrentmap.ConcurrentMap,
	v1 interface{},
	v2 interface{},
	err error,
) {
	cm, found := i.Env.Get(arg1)
	if !found {
		return cm, v1, v2, fmt.Errorf("Invalid variable %q", arg1)
	}
	cm.Lock()
	variablei, _ := cm.GetLockless(arg1)
	variable, _ := variablei.(*vrb.Variable)
	variableType := variable.Type
	cm.Unlock() // other functions require lock
	switch variableType {
	case vrb.FLOAT:
		v1, err = i.GetBigFloat64(arg1)
		if err != nil {
			return cm, v1, v2, err
		}
		v2, err = i.GetBigFloat64(arg2)
		if err != nil {
			return cm, v1, v2, err
		}
	case vrb.INT:
		v1, err = i.GetInt64(arg1)
		if err != nil {
			return cm, v1, v2, err
		}
		v2, err = i.GetInt64(arg2)
		if err != nil {
			return cm, v1, v2, err
		}
	default:
		return cm, v1, v2,
			fmt.Errorf("invalid variable type %q", variableType)
	}
	return cm, v1, v2, nil
}
