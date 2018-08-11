package instruction

import (
	"posam/interpreter"
	"posam/util/concurrentmap"
)

type Instructioner interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type Instruction struct {
	Env   *concurrentmap.ConcurrentMap
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

func (i *Instruction) ParseVariable(name string) (*interpreter.Variable, error) {
	v, found := i.Env.Get(name)
	if !found {
		newVariable := &interpreter.Variable{}
		newVariablei := i.Env.Set(name, newVariable)
		variable := newVariablei.(*interpreter.Variable)
		return variable, nil
	} else {
		variable := v.(*interpreter.Variable)
		return variable, nil
	}
}
