package instruction_test

import (
	"math/big"
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"synthesis/internal/interpreter/vrb"
	"testing"
)

func TestInstructionMultiplicationExecuteFloat64(t *testing.T) {
	i := instruction.InstructionMultiplication{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11.11")
	var2, _ := vrb.NewVariable("var2", "22.22")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "246.8642" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"246.8642",
			v.Value,
		)
	}

	i.Execute(var1.Name, "33.33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "8227.983786" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"8227.983786",
			v.Value,
		)
	}
}

func TestInstructionMultiplicationExecuteInt64(t *testing.T) {
	i := instruction.InstructionMultiplication{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11")
	var2, _ := vrb.NewVariable("var2", "22")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != 242 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			242,
			v.Value,
		)
	}

	i.Execute(var1.Name, "33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != 7986 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			7986,
			v.Value,
		)
	}
}
