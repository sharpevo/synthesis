package instruction_test

import (
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"testing"
)

func TestInstructionControlFlowGotoExecute(t *testing.T) {
	i := instruction.InstructionControlFlowGoto{}
	i.Env = interpreter.NewStack()
	i.Goto(3)
	variable, _ := i.Env.Get("SYS_CUR")
	actual := variable.Value.(int64)
	if actual != 3 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			3,
			actual,
		)
	}
}
