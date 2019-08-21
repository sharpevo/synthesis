package instruction_test

import (
	"fmt"
	"synthesis/instruction"
	"synthesis/interpreter"
	"testing"
)

func TestInstructionPrinterLoadCycle(t *testing.T) {
	i := instruction.InstructionPrinterLoadCycle{}
	i.Env = interpreter.NewStack()

	resp, err := i.Execute("var1", "tests/test.bin")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
	expected := 3
	if v, _ := i.Env.Get("var1"); v.Value.(int) != expected {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			expected,
			v.Value,
		)
	}
}
