package instruction_test

import (
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionVariableGet(t *testing.T) {
	i := instruction.InstructionVariableGet{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &vrb.Variable{Value: "11.22"})

	resp, _ := i.Execute("var1")
	if !strings.Contains(resp.(string), "11.22") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.22",
			resp,
		)
	}

	resp, _ = i.Execute("var2")
	if !strings.Contains(resp.(string), "= <nil>") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"",
			resp,
		)
	}
}
