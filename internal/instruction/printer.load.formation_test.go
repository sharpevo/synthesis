package instruction_test

import (
	"fmt"
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"testing"
)

func TestInstructionPrinterLoadFormation(t *testing.T) {
	i := instruction.InstructionPrinterLoadFormation{}
	i.Env = interpreter.NewStack()

	resp, err := i.Execute("var1", "tests/test.bin", "0")
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
