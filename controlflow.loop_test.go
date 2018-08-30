package instruction_test

import (
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"testing"
)

func TestInstructionControlFlowLoopExecute(t *testing.T) {
	i := instruction.InstructionControlFlowLoop{}
	i.Env = interpreter.NewStack()
	varCur, _ := i.Env.Get("SYS_CUR")
	actual := varCur.Value.(int64)
	if actual != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			0,
			actual,
		)
	}

	line := "5"
	count := 3

	varCount, _ := vrb.NewVariable("loopcount", "3")
	i.Env.Set(varCount)
	for index := range [3]int{} {
		resp, err := i.Execute(line, varCount.Name)
		t.Logf("%v", resp)
		if err != nil {
			t.Fatal(err)
		}
		expectedCount := count - index - 1
		if expectedCount > 1 && varCount.Value.(int64) != int64(expectedCount) {
			t.Errorf(
				"\nEXPECT: %v\nGET: %v\n",
				expectedCount,
				varCount.Value,
			)
		}

	}

}
