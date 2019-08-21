package instruction_test

import (
	"synthesis/instruction"
	"synthesis/interpreter"
	"testing"
)

func TestInstructionControlFlowErrGotoExecute(t *testing.T) {
	i := instruction.InstructionControlFlowErrGoto{}
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
	varErr, _ := i.Env.Get("SYS_ERR")
	varErr.Value = "err occured"
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
