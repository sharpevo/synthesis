package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"testing"
)

func TestInstructionControlFlowGotoGoto(t *testing.T) {
	i := instruction.InstructionControlFlowGoto{}
	i.Env = interpreter.NewStack()
	i.Goto(3)
	variable, _ := i.Env.Get("SYS_CUR")
	actual := variable.Value.(int)
	if actual != 3 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			3,
			actual,
		)
	}
}
