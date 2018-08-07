package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"testing"
)

func TestInstructionAdditionFloat64Execute(t *testing.T) {
	i := instruction.InstructionAdditionFloat64{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &interpreter.Variable{Value: "12.34"})
	i.Env.Set("var2", &interpreter.Variable{Value: "43.21"})

	i.Execute("var1", "var2")
	if v, _ := i.Env.Get("var1"); v.(*interpreter.Variable).Value != "55.55" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.(*interpreter.Variable).Value,
		)
	}

	i.Execute("var1", "33.33")
	if v, _ := i.Env.Get("var1"); v.(*interpreter.Variable).Value != "88.88" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"88.88",
			v.(*interpreter.Variable).Value,
		)
	}
}
