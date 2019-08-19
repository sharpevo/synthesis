package instruction_test

import (
	"math/big"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"testing"
)

func TestInstructionSubtractionExecuteFloat64(t *testing.T) {
	i := instruction.InstructionSubtraction{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "88.88")
	var2, _ := vrb.NewVariable("var2", "11.11")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "77.77" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"77.77",
			v.Value,
		)
	}

	i.Execute(var1.Name, "22.22")
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "55.55" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.Value,
		)
	}
}

func TestInstructionSubtractionExecuteInt64(t *testing.T) {
	i := instruction.InstructionSubtraction{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11")
	var2, _ := vrb.NewVariable("var2", "22")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != -11 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			-11,
			v.Value,
		)
	}

	i.Execute(var1.Name, "33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != -44 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			-44,
			v.Value,
		)
	}
}
