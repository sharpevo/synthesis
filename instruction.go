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
	variable, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, name)
		if err != nil {
			return variable, err
		}
		variable = i.Env.Set(newVariable)
	}
	return variable, nil
}
