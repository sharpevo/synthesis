package instruction_test

import (
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
	"testing"
)

func TestInstructionMultiplicationFloat64Execute(t *testing.T) {
	i := instruction.InstructionMultiplicationFloat64{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &vrb.Variable{Value: "11.11"})
	i.Env.Set("var2", &vrb.Variable{Value: "22.22"})

	i.Execute("var1", "var2")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "246.8642" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"246.8642",
			v.(*vrb.Variable).Value,
		)
	}

	i.Execute("var1", "33.33")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "8227.983786" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"8227.983786",
			v.(*vrb.Variable).Value,
		)
	}
}
