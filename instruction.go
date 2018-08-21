package instruction

import (
	"posam/interpreter"
	"posam/interpreter/vrb"
)

type Instructioner interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type Instruction struct {
	Env   *interpreter.Stack
	title string
	path  string
}

func (i Instruction) Title() string {
	return i.title
}

func (i *Instruction) SetTitle(title string) {
	i.title = title
}

func (i *Instruction) Execute(args ...string) (interface{}, error) {
	return "", nil
}

func (i *Instruction) ParseVariable(name string) (*vrb.Variable, error) {
	v, found := i.Env.Get(name)
	if !found {
		newVariable := &vrb.Variable{}
		newVariablei := i.Env.Set(name, newVariable)
		variable := newVariablei.(*vrb.Variable)
		return variable, nil
	} else {
		variable := v.(*vrb.Variable)
		return variable, nil
	}
}
