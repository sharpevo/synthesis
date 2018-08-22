package instruction_test

import (
	"fmt"
	"posam/instruction"
	"posam/interpreter"
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
	for index := range [3]int{} {
		_, err := i.Execute(line, fmt.Sprintf("%d", count))
		if err != nil {
			t.Fatal(err)
		}
		expectedCount := index + 1
		if i.Count() != expectedCount {
			t.Errorf(
				"\nEXPECT: %v\nGET: %v\n",
				expectedCount,
				i.Count(),
			)
		}

	}

}
