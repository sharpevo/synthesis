package instruction

import (
	"fmt"
	"posam/dao"
	"strings"
)

func init() {
	dao.InstructionMap.Set("GETVARBYINDEX", InstructionVariableGet{})
}

type InstructionVariableGetByIndex struct {
	Instruction
}

func (i *InstructionVariableGetByIndex) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
		return resp, fmt.Errorf("not enough arguments")
	}
	targetName := args[0]
	targetCM, err := i.ParseFloat64Variable(targetName)
	if err != nil {
		return resp, err
	}
	index, err := i.ParseInt(args[2])
	if err != nil {
		return resp, err
	}

	variableName := args[1]
	variableCM, found := i.Env.Get(variableName)
	if !found {
		resp = fmt.Sprintf("%s is not defined", variableName)
		return
	}
	isSameCM := targetCM == variableCM
	variableCM.Lock()
	defer variableCM.Unlock()
	variable, _ := i.GetVarLockless(variableCM, variableName)
	list := strings.Split(fmt.Sprintf("%s", variable.GetValue()), ",")
	if index > len(list)-1 || index < 0 {
		return resp, fmt.Errorf("invalid index %v from %v", index, list)
	}
	if !isSameCM {
		targetCM.Lock()
	}
	targetVar, _ := i.GetVarLockless(targetCM, targetName)
	targetVar.SetValue(list[index])
	resp = fmt.Sprintf("%s = %s (%s[%v])", targetName, list[index], variableName, index)
	if !isSameCM {
		targetCM.Unlock()
	}
	return resp, nil
}
