package instruction

import (
	"fmt"
	"posam/dao"
	"posam/interpreter/vrb"
)

func init() {
	dao.InstructionMap.Set("CMPVAR", InstructionVariableCompare{})
}

type InstructionVariableCompare struct {
	Instruction
}

func (i *InstructionVariableCompare) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name1 := args[0]
	name2 := args[1]
	cm1, err := i.ParseVariable(name1)
	if err != nil {
		return resp, err
	}
	cm2, err := i.ParseVariable(name2)
	if err != nil {
		return resp, err
	}
	isSameCM := cm1 == cm2
	cm1.Lock()
	if !isSameCM {
		cm2.Lock()
	}
	var1, _ := i.GetVarLockless(cm1, name1)
	var2, _ := i.GetVarLockless(cm2, name2)
	result, err := vrb.Compare(var1, var2)
	cm1.Unlock()
	if !isSameCM {
		cm2.Unlock()
	}
	if err != nil {
		return resp, err
	}
	cm, found := i.Env.Get("SYS_CMP")
	if !found {
		return resp, fmt.Errorf("failed to load variable CMP")
	}
	cm.Lock()
	v, _ := i.GetVarLockless(cm, "SYS_CMP")
	v.SetValue(result)
	cm.Unlock()
	return result.String(), nil
}
