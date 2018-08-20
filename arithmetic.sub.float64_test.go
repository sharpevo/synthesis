package instruction_test

import (
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
	"testing"
)

func TestInstructionSubtractionFloat64Execute(t *testing.T) {
	i := instruction.InstructionSubtractionFloat64{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &vrb.Variable{Value: "88.88"})
	i.Env.Set("var2", &vrb.Variable{Value: "11.11"})

	i.Execute("var1", "var2")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "77.77" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"77.77",
			v.(*vrb.Variable).Value,
		)
	}

	i.Execute("var1", "22.22")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "55.55" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.(*vrb.Variable).Value,
		)
	}
}
