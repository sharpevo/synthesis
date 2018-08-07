package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionSubtractionFloat64Execute(t *testing.T) {
	var err error
	i := instruction.InstructionSubtractionFloat64{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &interpreter.Variable{Value: "88.88"})
	i.Env.Set("var2", &interpreter.Variable{Value: "11.11"})

	// addition of float variable
	_, err = i.Execute("var1", "var2")
	if v, _ := i.Env.Get("var1"); v.(*interpreter.Variable).Value != "77.77" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"77.77",
			v.(*interpreter.Variable).Value,
		)
	}

	// addition of variable and float
	_, err = i.Execute("var1", "22.22")
	if v, _ := i.Env.Get("var1"); v.(*interpreter.Variable).Value != "55.55" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.(*interpreter.Variable).Value,
		)
	}

	// addition of invalid variable
	_, err = i.Execute("invalid", "33.33")
	if !strings.Contains(err.Error(), "Invalid variable") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Invalid variable",
			err.Error(),
		)
	}

	// addition of variable that is not Variable pointer
	i.Env.Set("var3", interpreter.Variable{Value: "not variable pointer"})
	_, err = i.Execute("var3", "33.33")
	if !strings.Contains(err.Error(), "Invalid type of variable") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Invalid variable",
			err.Error(),
		)
	}
}
