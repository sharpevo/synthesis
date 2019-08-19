package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"strings"
	"testing"
)

func TestInstructionVariableGet(t *testing.T) {
	i := instruction.InstructionVariableGet{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11.22")
	i.Env.Set(var1)

	resp, _ := i.Execute(var1.Name)
	if !strings.Contains(resp.(string), "11.22") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.22",
			resp,
		)
	}

	resp, _ = i.Execute("var2")
	if !strings.Contains(resp.(string), "not defined") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"",
			resp,
		)
	}
}
