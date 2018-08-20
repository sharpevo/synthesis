package instruction

import (
	"fmt"
	"posam/interpreter/vrb"
	"strings"
)

type InstructionVariableSet struct {
	Instruction
}

func (i *InstructionVariableSet) Execute(args ...string) (resp interface{}, err error) {
	if len(args) <= 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	name := args[0]
	value := strings.Join(args[1:], " ")
	v, found := i.Env.Get(name)
	if !found {
		variable, err := vrb.NewVariable(name, value)
		if err != nil {
			return resp, err
		}
		resp = i.Env.Set(name, variable)
		return fmt.Sprintf(
			"variable %q is set to \"%v\"",
			name,
			value,
		), nil
	} else {
		variable := v.(*vrb.Variable)
		variable.Value, variable.Type, _ = vrb.ParseValue(value)
		return resp, nil
	}
}
