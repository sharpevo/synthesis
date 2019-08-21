package instruction

import (
	"fmt"
	"log"
	"synthesis/interpreter/vrb"
	"synthesis/util/blockingqueue"
	"synthesis/util/concurrentmap"
	"sync"
)

type Stack struct {
	// TODO: mutex
	lock  sync.RWMutex
	cmaps *blockingqueue.BlockingQueue
}

func NewStack(stackList ...*Stack) *Stack {
	newMap := concurrentmap.NewConcurrentMap()
	preservedVariableList := vrb.NewPreservedVariables()
	for _, v := range preservedVariableList {
		newMap.Set(v.Name, v)
	}
	newStack := &Stack{
		cmaps: blockingqueue.NewBlockingQueue(),
	}
	newStack.cmaps.Append(newMap)
	if len(stackList) != 0 {
		for _, stack := range stackList {
			for item := range stack.cmaps.Iter() {
				newStack.cmaps.Append(item.Value)
			}
		}
	}
	return newStack
}

func (s *Stack) Len() int {
	return s.cmaps.Length()
}

func (s *Stack) GetVariable(name string) (variable *vrb.Variable, found bool) {
	for item := range s.cmaps.Iter() {
		if found {
			continue
		}
		cmapi := item.Value
		cmap := cmapi.(*concurrentmap.ConcurrentMap)
		cmap.Lock()
		if v, found := cmap.GetLockless(name); found {
			variable = v.(*vrb.Variable)
			fmt.Printf("reading stack %#v: %v\n", name, variable)
		}
		cmap.Unlock()
	}
	return variable, found
}

func (s *Stack) Get(name string) (cmap *concurrentmap.ConcurrentMap, found bool) {
	for item := range s.cmaps.Iter() {
		if found {
			continue
		}
		cmapi := item.Value
		_cmap := cmapi.(*concurrentmap.ConcurrentMap)
		_cmap.Lock()
		if _, found = _cmap.GetLockless(name); found {
			cmap = _cmap
			found = true
		}
		_cmap.Unlock()
	}
	return cmap, found
}

func (s *Stack) Set(variable *vrb.Variable) *vrb.Variable {
	// TODO: global variable creation
	cmapi, _ := s.cmaps.Get(0)
	cmap := cmapi.(*concurrentmap.ConcurrentMap)
	result := cmap.Set(variable.Name, variable)
	log.Printf("Set stack: %s = %v\n", variable.Name, variable.GetValue())
	return result.(*vrb.Variable)
}

func (s *Stack) Push(cmap *concurrentmap.ConcurrentMap) error {
	for item := range s.cmaps.Iter() {
		_cmapi := item.Value
		_cmap := _cmapi.(*concurrentmap.ConcurrentMap)
		if _cmap == cmap {
			return fmt.Errorf("item duplicated in stack")
		}
	}
	s.cmaps.Append(cmap)
	return nil
}

func (s *Stack) Lock() {
	s.lock.Lock()
}

func (s *Stack) Unlock() {
	s.lock.Unlock()
}
