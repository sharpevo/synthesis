package instruction

import (
	"fmt"
	"synthesis/internal/dao"
	"synthesis/internal/interpreter/vrb"
	"strings"
)

func init() {
	dao.InstructionMap.Set("SETVAR", InstructionVariableSet{})
}

type InstructionVariableSet struct {
	Instruction
}

func (i *InstructionVariableSet) Execute(args ...string) (resp interface{}, err error) {
	if len(args) <= 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name := args[0]
	value := strings.Join(args[1:], " ")
	cm, found := i.Env.Get(name)
	if !found {
		variable, err := vrb.NewVariable(name, value)
		if err != nil {
			return resp, err
		}
		i.Env.Set(variable)
	} else {
		cm.Lock()
		variable, _ := i.GetVarLockless(cm, name)
		v, t, _ := vrb.ParseValue(value)
		variable.Type = t
		variable.SetValue(v)
		cm.Unlock()
	}
	return fmt.Sprintf(
		"variable %q is set to \"%v\"",
		name,
		value,
	), nil
}
