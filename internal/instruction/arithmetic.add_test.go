package instruction_test

import (
	"math/big"
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"synthesis/internal/interpreter/vrb"
	"testing"
)

func TestInstructionAdditionExecuteFloat64(t *testing.T) {
	i := instruction.InstructionAddition{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "12.34")
	var2, _ := vrb.NewVariable("var2", "43.21")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "55.55" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.Value,
		)
	}

	i.Execute(var1.Name, "33.33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "88.88" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"88.88",
			v.Value,
		)
	}
}

func TestInstructionAdditionExecuteInt64(t *testing.T) {
	i := instruction.InstructionAddition{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11")
	var2, _ := vrb.NewVariable("var2", "22")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != 33 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			33,
			v.Value,
		)
	}

	i.Execute(var1.Name, "33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != 66 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			66,
			v.Value,
		)
	}
}
