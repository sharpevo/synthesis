package instruction_test

import (
	"fmt"
	"math/big"
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestParseObjects(t *testing.T) {
	var variable *interpreter.Variable
	var v1 *big.Float
	var v2 *big.Float
	var err error

	i := instruction.InstructionArithmetic{}
	i.Env = concurrentmap.NewConcurrentMap()
	i.Env.Set("var1", &interpreter.Variable{Value: "11.11"})
	i.Env.Set("var2", &interpreter.Variable{Value: "22.22"})
	i.Env.Set("var3", interpreter.Variable{Value: "not variable pointer"})

	variable, v1, v2, err = i.ParseObjects("var1", "var2")
	if v1.Cmp(big.NewFloat(11.11)) != 0 ||
		v2.Cmp(big.NewFloat(22.22)) != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.11 and 22.22",
			fmt.Sprintf("%s and %s", v1.String(), v2.String()),
		)
	}

	variable, v1, v2, err = i.ParseObjects("var1", "33.33")
	if v1.Cmp(big.NewFloat(11.11)) != 0 ||
		v2.Cmp(big.NewFloat(33.33)) != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.11 and 33.33",
			fmt.Sprintf("%s and %s", v1.String(), v2.String()),
		)
	}

	variable, v1, v2, err = i.ParseObjects("var1", "var3")
	if !strings.Contains(err.Error(), "invalid float") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"invalid float",
			err.Error(),
		)
	}

	variable, v1, v2, err = i.ParseObjects("var3", "33.33")
	if !strings.Contains(err.Error(), "Invalid type of variable") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Invalid type of variable",
			err.Error(),
		)
	}

	variable, v1, v2, err = i.ParseObjects("invalid", "33.33")
	if !strings.Contains(err.Error(), "Invalid variable") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Invalid variable",
			err.Error(),
		)
	}
	fmt.Println(variable)

}

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
