package instruction_test

import (
	"fmt"
	"math/big"
	"strings"
	"synthesis/instruction"
	"synthesis/interpreter"
	"synthesis/interpreter/vrb"
	"testing"
)

func TestParseObjects(t *testing.T) {
	var variable *vrb.Variable
	var v1 *big.Float
	var v2 *big.Float
	var err error

	i := instruction.InstructionArithmetic{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11.11")
	var2, _ := vrb.NewVariable("var2", "22.22")
	var3, _ := vrb.NewVariable("var3", "string")
	i.Env.Set(var1)
	i.Env.Set(var2)
	i.Env.Set(var3)

	variable, v1, v2, err = i.ParseObjects(var1.Name, var2.Name)
	if v1.Cmp(big.NewFloat(11.11)) != 0 ||
		v2.Cmp(big.NewFloat(22.22)) != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.11 and 22.22",
			fmt.Sprintf("%s and %s", v1.String(), v2.String()),
		)
	}

	variable, v1, v2, err = i.ParseObjects(var1.Name, "33.33")
	if v1.Cmp(big.NewFloat(11.11)) != 0 ||
		v2.Cmp(big.NewFloat(33.33)) != 0 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11.11 and 33.33",
			fmt.Sprintf("%s and %s", v1.String(), v2.String()),
		)
	}

	variable, v1, v2, err = i.ParseObjects(var1.Name, var3.Name)
	if !strings.Contains(err.Error(), "syntax error scanning number") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"invalid float",
			err.Error(),
		)
	}

	variable, v1, v2, err = i.ParseObjects(var3.Name, "33.33")
	if !strings.Contains(err.Error(), "syntax error scanning number") {
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

func TestParseOperandsInt(t *testing.T) {
	var variable *vrb.Variable
	var v1 interface{}
	var v2 interface{}
	var err error

	i := instruction.InstructionArithmetic{}
	i.Env = interpreter.NewStack()
	var1, _ := vrb.NewVariable("var1", "11")
	var2, _ := vrb.NewVariable("var2", "22")
	var3, _ := vrb.NewVariable("var3", "string")
	i.Env.Set(var1)
	i.Env.Set(var2)
	i.Env.Set(var3)

	variable, v1, v2, err = i.ParseOperands(var1.Name, var2.Name)
	if v1.(int64) != 11 || v2.(int64) != 22 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11 and 22",
			fmt.Sprintf("%v and %v", v1, v2),
		)
	}

	variable, v1, v2, err = i.ParseOperands(var1.Name, "33")
	if v1.(int64) != 11 || v2.(int64) != 33 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"11 and 33",
			fmt.Sprintf("%v and %v", v1, v2),
		)
	}

	variable, v1, v2, err = i.ParseOperands(var1.Name, var3.Name)
	if !strings.Contains(err.Error(), "invalid syntax") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"invalid syntax",
			err.Error(),
		)
	}

	variable, v1, v2, err = i.ParseOperands(var3.Name, "33")
	if !strings.Contains(err.Error(), "invalid variable type") {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"Invalid type of variable",
			err.Error(),
		)
	}

	variable, v1, v2, err = i.ParseOperands("invalid", "33")
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
	arithmetic.Env = interpreter.NewStack()

	nonFloatVar, _ := vrb.NewVariable("nonfloat", "none float")
	v, err = arithmetic.GetBigFloat64(nonFloatVar)
	if err == nil {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"not match error",
			err,
		)
	}

	variable, _ := vrb.NewVariable("testname", "55.55")
	arithmetic.Env.Set(variable)
	varVariable, _ := arithmetic.Env.Get(variable.Name)

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

func TestGetInt64(t *testing.T) {
	var v int64
	var err error

	arithmetic := instruction.InstructionArithmetic{}
	arithmetic.Env = interpreter.NewStack()

	nonIntVar, _ := vrb.NewVariable("nonint", "none int")
	v, err = arithmetic.GetInt64(nonIntVar)
	if err == nil {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			"not match error",
			err,
		)
	}

	variable, _ := vrb.NewVariable("testname", "10")
	arithmetic.Env.Set(variable)
	varVariable, _ := arithmetic.Env.Get(variable.Name)

	v, err = arithmetic.GetInt64(varVariable)
	if err != nil || v != 10 {
		t.Log(err)
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			10,
			v,
		)
	}

	varIntString := "11"
	v, err = arithmetic.GetInt64(varIntString)
	if err != nil || v != 11 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			12,
			v,
		)
	}

	varInt64 := int64(14)
	v, err = arithmetic.GetInt64(varInt64)
	if err != nil || v != 14 {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			14,
			v,
		)
	}

	varInvalidString := "test"
	v, err = arithmetic.GetInt64(varInvalidString)
	if !strings.Contains(err.Error(), "invalid syntax") ||
		v != 0 {
		t.Logf(err.Error())
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			0,
			v,
		)
	}
}
