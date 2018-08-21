package interpreter_test

import (
	"posam/interpreter"
	"posam/interpreter/vrb"
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
	for k, _ := range vrb.PreservedNames {
		_, ok := s1.Get(k)
		if !ok {
			t.Errorf(
				"\nEXPECT: %v: %q\nGET:%v\n",
				"init stack with preserved variable", k,
				"not found",
			)
		}
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
	variable, err := vrb.NewVariable("var1", "test")
	if err != nil {
		t.Fatal(err)
	}
	stack.Set(variable.Name, variable)
	if v, found := stack.Get("var1"); !found {
		t.Errorf(
			"\nEXPECT: %v\nGET:%v\n",
			v,
			"not found",
		)
	}
}
