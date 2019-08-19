package instruction_test

import (
	"fmt"
	"math/big"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"testing"
)

func TestInstructionModuloExecuteFloat64(t *testing.T) {
	i := instruction.InstructionModulo{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "28.02")
	var2, _ := vrb.NewVariable("var2", "2.0")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "0" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0",
			v.Value,
		)
	}

	var3, _ := vrb.NewVariable("var3", "20.1")
	i.Env.Set(var3)
	i.Execute(var3.Name, "5")
	fmt.Println(var3.Value.(*big.Float).String())
	if v, _ := i.Env.Get(var3.Name); v.Value.(*big.Float).String() != "0" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0",
			v.Value,
		)
	}
}

func TestInstructionModuloExecuteInt64(t *testing.T) {
	i := instruction.InstructionModulo{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "33")
	var2, _ := vrb.NewVariable("var2", "3")
	i.Env.Set(var1)
	i.Env.Set(var2)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(int64) != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			0,
			v.Value,
		)
	}

	var3, _ := vrb.NewVariable("var3", "20")
	i.Env.Set(var3)
	i.Execute(var3.Name, "3")
	if v, _ := i.Env.Get(var3.Name); v.Value.(int64) != 2 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			0,
			v.Value,
		)
	}
}
