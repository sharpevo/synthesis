package interpreter

import (
	"fmt"
	"posam/util/concurrentmap"
)

type Stack struct {
	// TODO: mutex
	cmaps []*concurrentmap.ConcurrentMap
}

func NewStack(stackList ...*Stack) *Stack {
	newMap := concurrentmap.NewConcurrentMap()
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

func (s *Stack) Get(name string) (interface{}, bool) {
	for _, cmap := range s.cmaps {
		if v, found := cmap.Get(name); found {
			return v, found
		}
	}
	return nil, false
}

func (s *Stack) Set(name string, value interface{}) interface{} {
	// TODO: global variable creation
	cmap := s.cmaps[0]
	result := cmap.Set(name, value)
	return result
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
