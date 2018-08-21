package interpreter_test

import (
	"posam/interpreter"
	"posam/util/concurrentmap"
	"testing"
)

func TestNewStackAndLenAndPush(t *testing.T) {
	s1 := interpreter.NewStack()
	if s1.Len() != 1 {
		t.Errorf(
			"\nEXPECT: %v\nGET:%v\n",
			1,
			s1.Len(),
		)
	}
	s1.Push(concurrentmap.NewConcurrentMap())
	s2 := interpreter.NewStack(s1)
	if s2.Len() != 3 {
		t.Errorf(
			"\nEXPECT: %v\nGET:%v\n",
			3,
			s2.Len(),
		)
	}
	t.Logf("%#v", s2)
}

func TestStackGetAndSet(t *testing.T) {
	stack := interpreter.NewStack()
	if v, found := stack.Get("var1"); found {
		t.Errorf(
			"\nEXPECT: %v\nGET:%v\n",
			"not found",
			v,
		)
	}
	stack.Set("var1", "test")
	if v, found := stack.Get("var1"); !found {
		t.Errorf(
			"\nEXPECT: %v\nGET:%v\n",
			v,
			"not found",
		)
	}
}
