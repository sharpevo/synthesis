package instruction

import (
	"fmt"
	"synthesis/internal/dao"
)

func init() {
	dao.InstructionMap.Set("GETVAR", InstructionVariableGet{})
}

type InstructionVariableGet struct {
	Instruction
}

func (i *InstructionVariableGet) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name := args[0]
	cm, found := i.Env.Get(name)
	if !found {
		resp = fmt.Sprintf("%s is not defined", name)
		return
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, name)
	resp = fmt.Sprintf("%v %s = %v", variable.Type, name, variable.GetValue())
	cm.Unlock()
	return resp, nil
}
