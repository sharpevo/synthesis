package instruction_test

import (
	"math/big"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"strings"
	"testing"
)

func TestInstructionDivisionFloat64Execute(t *testing.T) {
	i := instruction.InstructionDivisionFloat64{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11.11")
	var2, _ := vrb.NewVariable("var2", "22.22")
	var3, _ := vrb.NewVariable("var3", "0.0")
	i.Env.Set(var1)
	i.Env.Set(var2)
	i.Env.Set(var3)

	i.Execute(var1.Name, var2.Name)
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "0.5" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0.50",
			v.Value,
		)
	}

	i.Execute(var1.Name, "33.33")
	if v, _ := i.Env.Get(var1.Name); v.Value.(*big.Float).String() != "0.01500150015" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0.01500150015",
			v.Value,
		)
	}

	_, err := i.Execute(var1.Name, "0")
	if !strings.Contains(err.Error(), "quotition") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"inf quotition",
			err.Error(),
		)
	}

	i.Execute(var3.Name, "33.33")
	if v, _ := i.Env.Get(var3.Name); v.Value.(*big.Float).String() != "0" {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			"0",
			v.Value,
		)
	}

}
