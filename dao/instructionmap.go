package dao

import (
	"posam/util/concurrentmap"
	"reflect"
	"sync"
)

type InstructionMapt struct {
	lock sync.Mutex
	cmap *concurrentmap.ConcurrentMap
}

func NewInstructionMap() *InstructionMapt {
	return &InstructionMapt{
		cmap: concurrentmap.NewConcurrentMap(),
	}
}

func (m *InstructionMapt) Set(k string, v interface{}) {
	value, ok := v.(reflect.Type)
	if ok {
		m.cmap.Set(k, value)
	} else {
		m.cmap.Set(k, reflect.TypeOf(v))
	}
}

func (m *InstructionMapt) Lock() {
	m.cmap.Lock()
}

func (m *InstructionMapt) Unlock() {
	m.cmap.Unlock()
}

func (m *InstructionMapt) Iter() <-chan concurrentmap.Item {
	return m.cmap.Iter()
}