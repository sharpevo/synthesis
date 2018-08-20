package instruction

import (
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
)

type Instructioner interface {
	Execute(args ...string) (interface{}, error)
	SetEnv(*concurrentmap.ConcurrentMap)
	GetEnv() *concurrentmap.ConcurrentMap
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

// TODO: Env with tons of reflect relavent jobs
func (i *Instruction) GetEnv() *concurrentmap.ConcurrentMap {
	if i.Env == nil {
		i.initEnv()
	}
	return i.Env
}

func (i *Instruction) SetEnv(env *concurrentmap.ConcurrentMap) {
	i.Env = env
}

func (i *Instruction) initEnv() {
	i.Env = concurrentmap.NewConcurrentMap()
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
