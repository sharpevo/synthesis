package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"testing"
)

func TestParseVariable(t *testing.T) {
	i := instruction.Instruction{}

	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "var1 test")
	i.Env.Set(var1)

	v1, _ := i.ParseVariable(var1.Name)
	if v1.Value != "var1 test" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Test",
			v1.Value,
		)
	}
	v2, _ := i.ParseVariable("var2")
	v2.Value = "test var2"

	var2, _ := i.Env.Get("var2")
	if var2.Value != "test var2" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"test var2",
			var2.Value,
		)
	}
}
