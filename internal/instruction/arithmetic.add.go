package instruction

import (
	"fmt"
	"math/big"
	"synthesis/internal/dao"
	"synthesis/internal/interpreter/vrb"
)

func init() {
	dao.InstructionMap.Set("ADD", InstructionAddition{})
}

type InstructionAddition struct {
	InstructionArithmetic
}

func (i *InstructionAddition) Execute(args ...string) (resp interface{}, err error) {
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
		v1.Add(v1, v2)
		variable.SetValue(v1)
		return v1.String(), nil
	case vrb.INT:
		v1 := v1v.(int64)
		v2 := v2v.(int64)
		v1 += v2
		variable.SetValue(v1)
		return fmt.Sprintf("%v", v1), nil
	default:
		return "", fmt.Errorf("invalid variable type")
	}
	return resp, err
}
