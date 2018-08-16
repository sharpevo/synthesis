package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"testing"
)

func TestParseVariable(t *testing.T) {
	i := instruction.Instruction{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &interpreter.Variable{Value: "var1 test"})
	v1, _ := i.ParseVariable("var1")
	if v1.Value != "var1 test" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Test",
			v1.Value,
		)
	}
	v2, _ := i.ParseVariable("var2")
	v2.Value = "test var2"

	v2i, _ := i.Env.Get("var2")
	var2 := v2i.(*interpreter.Variable)
	if var2.Value != "test var2" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"test var2",
			var2.Value,
		)
	}
}
