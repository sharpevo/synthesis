package instruction

import (
	"fmt"
	"log"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
	"sync"
)

type Stack struct {
	// TODO: mutex
	lock  sync.RWMutex
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

func (s *Stack) GetVariable(name string) (variable *vrb.Variable, found bool) {
	for _, cmap := range s.cmaps {
		cmap.Lock()
		if v, found := cmap.GetLockless(name); found {
			variable := v.(*vrb.Variable)
			fmt.Printf("reading stack %#v: %v\n", name, variable)
			cmap.Unlock()
			return variable, found
		}
		cmap.Unlock()
	}
	return variable, found
}

func (s *Stack) Get(name string) (cmap *concurrentmap.ConcurrentMap, found bool) {
	for _, cmap := range s.cmaps {
		cmap.Lock()
		if _, found = cmap.GetLockless(name); found {
			cmap.Unlock()
			return cmap, found
		}
		cmap.Unlock()
	}
	return cmap, found
}

func (s *Stack) Set(variable *vrb.Variable) *vrb.Variable {
	// TODO: global variable creation
	cmap := s.cmaps[0]
	result := cmap.Set(variable.Name, variable)
	log.Printf("Set stack: %s = %v\n", variable.Name, variable.GetValue())
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

func (s *Stack) Lock() {
	s.lock.Lock()
}

func (s *Stack) Unlock() {
	s.lock.Unlock()
}
