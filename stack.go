package interpreter

import (
	"fmt"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
)

type Stack struct {
	// TODO: mutex
	cmaps []*concurrentmap.ConcurrentMap
}

func NewStack(stackList ...*Stack) *Stack {
	newMap := concurrentmap.NewConcurrentMap()
	preservedVariableList := vrb.NewPreservedVariables()
	for _, v := range preservedVariableList {
		newMap.Set(v.Name, v)
	}
	newStack := &Stack{
		cmaps: []*concurrentmap.ConcurrentMap{newMap},
	}
	if len(stackList) != 0 {
		for _, stack := range stackList {
			newStack.cmaps = append(newStack.cmaps, stack.cmaps...)
		}
	}
	return newStack
}

func (s *Stack) Len() int {
	return len(s.cmaps)
}

func (s *Stack) Get(name string) (variable *vrb.Variable, found bool) {
	for _, cmap := range s.cmaps {
		if v, found := cmap.Get(name); found {
			return v.(*vrb.Variable), found
		}
	}
	return
}

func (s *Stack) Set(variable *vrb.Variable) *vrb.Variable {
	// TODO: global variable creation
	cmap := s.cmaps[0]
	result := cmap.Set(variable.Name, variable)
	return result.(*vrb.Variable)
}

func (s *Stack) Push(cmap *concurrentmap.ConcurrentMap) error {
	for _, _cmap := range s.cmaps {
		if _cmap == cmap {
			return fmt.Errorf("item duplicated in stack")
		}
	}
	s.cmaps = append(s.cmaps, cmap)
	return nil
}
