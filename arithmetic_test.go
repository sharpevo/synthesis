package instruction_test

import (
	"math/big"
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"testing"
)

func TestGetBigFloat64(t *testing.T) {
	var v *big.Float
	var err error

	arithmetic := instruction.InstructionArithmetic{}
	arithmetic.Env = concurrentmap.NewConcurrentMap()
	variable := interpreter.Variable{
		Value: "55.55",
	}
	v, err = arithmetic.GetBigFloat64(variable)
	if err == nil {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"not match error",
			err,
		)
	}

	arithmetic.Env.Set("testname", &variable)
	varVariable, _ := arithmetic.Env.Get("testname")

	v, err = arithmetic.GetBigFloat64(varVariable)
	if err != nil || v.String() != "55.55" {
		t.Log(err)
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"55.55",
			v.String(),
		)
	}

	varFloatString := "12.34"
	v, err = arithmetic.GetBigFloat64(varFloatString)
	if err != nil || v.String() != "12.34" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"12.34",
			v.String(),
		)
	}

	varBigFloat := big.NewFloat(43.21)
	v, err = arithmetic.GetBigFloat64(varBigFloat)
	if err != nil || v.String() != "43.21" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"43.21",
			v.String(),
		)
	}

	varInvalidString := "test"
	v, err = arithmetic.GetBigFloat64(varInvalidString)
	if err.Error() != "syntax error scanning number" ||
		v.String() != "0" {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"0",
			v.String(),
		)
	}
}
