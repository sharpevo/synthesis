package instruction

import (
	"fmt"
	"math/big"
	"posam/dao"
	"posam/interpreter/vrb"
)

func init() {
	dao.InstructionMap.Set("MOD", InstructionModulo{})
}

type InstructionModulo struct {
	InstructionArithmetic
}

func (i *InstructionModulo) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	cm, v1v, v2v, err := i.ParseOperands(args[0], args[1])
	if err != nil {
		return resp, err
	}
	cm.Lock()
	defer cm.Unlock()
	variablei, _ := cm.GetLockless(args[0])
	variable, _ := variablei.(*vrb.Variable)
	switch variable.Type {
	case vrb.FLOAT:
		v1 := v1v.(*big.Float)
		v2 := v2v.(*big.Float)
		if v2.String() == "0" {
			return resp, fmt.Errorf("modulo by zero")
		}
		v1Int, _ := v1.Int(nil)
		v2Int, _ := v2.Int(nil)
		v1Int.Mod(v1Int, v2Int)
		v1, _, err := big.ParseFloat(v1Int.String(), 10, 53, big.ToNearestEven)
		if err != nil {
			return resp, err
		}
		variable.SetValue(v1)
		return v1.String(), nil
	case vrb.INT:
		v1 := v1v.(int64)
		v2 := v2v.(int64)
		v1 %= v2
		variable.SetValue(v1)
		return fmt.Sprintf("%v", v1), nil
	default:
		return "", fmt.Errorf("invalid variable type")
	}
	return resp, err
}
