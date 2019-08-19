package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"testing"
)

func TestInstructionControlFlowGreaterThanGotoExecute(t *testing.T) {
	i := instruction.InstructionControlFlowGreaterThanGoto{}
	i.Env = interpreter.NewStack()
	i.Execute("5")
	varCur, _ := i.Env.Get("SYS_CUR")
	actual := varCur.Value.(int64)
	if actual != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			0,
			actual,
		)
	}
	varErr, _ := i.Env.Get("SYS_CMP")
	varErr.Value = vrb.GREATER
	i.Execute("10")
	actual = varCur.Value.(int64)
	if actual != 10 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			10,
			actual,
		)
	}
}
