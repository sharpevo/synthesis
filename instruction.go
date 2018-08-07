package instruction

import (
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
