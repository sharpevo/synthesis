package instruction_test

import (
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionDivisionFloat64Execute(t *testing.T) {
	i := instruction.InstructionDivisionFloat64{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &vrb.Variable{Value: "11.11"})
	i.Env.Set("var2", &vrb.Variable{Value: "22.22"})
	i.Env.Set("var3", &vrb.Variable{Value: "0"})

	i.Execute("var1", "var2")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "0.5" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0.50",
			v.(*vrb.Variable).Value,
		)
	}

	i.Execute("var1", "33.33")
	if v, _ := i.Env.Get("var1"); v.(*vrb.Variable).Value != "0.01500150015" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0.01500150015",
			v.(*vrb.Variable).Value,
		)
	}

	_, err := i.Execute("var1", "0")
	if !strings.Contains(err.Error(), "quotition") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"inf quotition",
			err.Error(),
		)
	}

	i.Execute("var3", "33.33")
	if v, _ := i.Env.Get("var3"); v.(*vrb.Variable).Value != "0" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0",
			v.(*vrb.Variable).Value,
		)
	}

}
