package instruction_test

import (
	"synthesis/instruction"
	"synthesis/interpreter"
	"synthesis/interpreter/vrb"
	"testing"
)

func TestInstructionControlFlowLessThanGotoExecute(t *testing.T) {
	i := instruction.InstructionControlFlowLessThanGoto{}
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
	varErr.Value = vrb.LESS
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
