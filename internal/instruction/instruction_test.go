package instruction_test

import (
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"synthesis/internal/interpreter/vrb"
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

func TestInstructionIssueError(t *testing.T) {
	i := instruction.Instruction{}
	i.Env = interpreter.NewStack()
	message := "test message"
	i.IssueError(message)
	varErr, _ := i.Env.Get("SYS_ERR")
	actual := varErr.Value.(string)
	if actual != message {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			message,
			actual,
		)
	}
}
