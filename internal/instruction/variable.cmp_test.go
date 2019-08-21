package instruction_test

import (
	"fmt"
	"synthesis/instruction"
	"synthesis/interpreter"
	"synthesis/interpreter/vrb"
	"strings"
	"testing"
)

func TestInstructionVariableCompareWithSampeType(t *testing.T) {
	i := instruction.InstructionVariableCompare{}
	i.Env = interpreter.NewStack()
	testList := []struct {
		var1     string
		var2     string
		expected vrb.ComparisonType
	}{
		// string
		{
			var1:     "string1",
			var2:     "string1",
			expected: vrb.EQUAL,
		},
		{
			var1:     "string2",
			var2:     "string3",
			expected: vrb.UNEQUAL,
		},
		// int
		{
			var1:     "1",
			var2:     "1",
			expected: vrb.EQUAL,
		},
		{
			var1:     "2",
			var2:     "3",
			expected: vrb.LESS,
		},
		{
			var1:     "5",
			var2:     "4",
			expected: vrb.GREATER,
		},
		// float
		{
			var1:     "10.24",
			var2:     "10.24",
			expected: vrb.EQUAL,
		},
		{
			var1:     "10.24",
			var2:     "40.96",
			expected: vrb.LESS,
		},
		{
			var1:     "81.92",
			var2:     "51.20",
			expected: vrb.GREATER,
		},
	}

	for k, test := range testList {
		var1Name := fmt.Sprintf("var%d-1", k)
		var2Name := fmt.Sprintf("var%d-2", k)
		var1, _ := vrb.NewVariable(var1Name, test.var1)
		var2, _ := vrb.NewVariable(var2Name, test.var2)
		i.Env.Set(var1)
		i.Env.Set(var2)
		i.Execute(var1Name, var2Name)
		resultVariable, _ := i.Env.Get("SYS_CMP")
		if resultVariable.Value != test.expected {
			t.Errorf(
				"\nEXPECT: %v\nGET:%v\n",
				test.expected,
				resultVariable.Value,
			)
		}
	}
}

func TestInstructionVariableCompareWithErrors(t *testing.T) {
	i := instruction.InstructionVariableCompare{}
	i.Env = interpreter.NewStack()
	testList := []struct {
		var1     string
		var2     string
		expected string
	}{
		{
			var1:     "string1",
			var2:     "0",
			expected: "mismatched",
		},
	}

	for k, test := range testList {
		var1Name := fmt.Sprintf("var%d-1", k)
		var2Name := fmt.Sprintf("var%d-2", k)
		var1, _ := vrb.NewVariable(var1Name, test.var1)
		var2, _ := vrb.NewVariable(var2Name, test.var2)
		i.Env.Set(var1)
		i.Env.Set(var2)
		_, err := i.Execute(var1Name, var2Name)
		if err == nil {
			t.Errorf("expected error not happened")
		} else {
			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf(
					"\nEXPECT: %v\nGET:%v\n",
					"mismatched error",
					err.Error(),
				)
			}
		}
	}
}
