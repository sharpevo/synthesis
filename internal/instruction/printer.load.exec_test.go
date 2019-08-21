package instruction_test

import (
	"fmt"
	"synthesis/instruction"
	"testing"
)

func TestInstructionPrinterLoadExec(t *testing.T) {
	i := instruction.InstructionPrinterLoadExec{}

	resp, err := i.Execute("tests/test.bin", "0", "1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
	//expected := 3
	//actual := resp.(int)
	//if actual != expected {
	//t.Errorf(
	//"\nEXPECT: %v\nGET: %v\n",
	//expected,
	//actual,
	//)
	//}
}
