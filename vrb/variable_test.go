package vrb_test

import (
	"fmt"
	"posam/interpreter/vrb"
	"testing"
)

func TestNewVariable(t *testing.T) {
	testList := []struct {
		name         string
		value        interface{}
		expectedType vrb.VariableType
	}{
		{
			name:         "var1",
			value:        "string value",
			expectedType: vrb.STRING,
		},
		{
			name:         "var2",
			value:        "1024",
			expectedType: vrb.INT,
		},
		{
			name:         "var3",
			value:        "1024.0",
			expectedType: vrb.FLOAT,
		},
		{
			name:         "var1",
			value:        "10.24",
			expectedType: vrb.FLOAT,
		},
	}

	for _, test := range testList {
		v, err := vrb.NewVariable(test.name, test.value.(string))
		if err != nil {
			t.Fatal(err)
		}
		actual := v.Type
		expected := test.expectedType
		if actual != expected {
			t.Errorf(
				"\nEXPECT: %v\n GET: %v\n",
				expected,
				actual,
			)
		}
	}
}

func TestCompare(t *testing.T) {
	testList := []struct {
		var1   string
		var2   string
		result vrb.ComparisonType
	}{
		// string
		{
			var1:   "string1",
			var2:   "string2",
			result: vrb.UNEQUAL,
		},
		{
			var1:   "string3",
			var2:   "string3",
			result: vrb.EQUAL,
		},

		// int
		{
			var1:   "1",
			var2:   "2",
			result: vrb.LESS,
		},
		{
			var1:   "3",
			var2:   "3",
			result: vrb.EQUAL,
		},
		{
			var1:   "5",
			var2:   "4",
			result: vrb.GREATER,
		},

		// float
		{
			var1:   "1.1",
			var2:   "1.2",
			result: vrb.LESS,
		},
		{
			var1:   "1.3",
			var2:   "1.30",
			result: vrb.EQUAL,
		},
		{
			var1:   "1.5",
			var2:   "1.4",
			result: vrb.GREATER,
		},

		// exceptions
		{
			var1:   "string5",
			var2:   "6",
			result: vrb.UNKNOWN,
		},
		{
			var1:   "7",
			var2:   "8.0",
			result: vrb.UNKNOWN,
		},
	}

	for k, test := range testList {
		t.Logf("#%d", k)
		var1, _ := vrb.NewVariable(
			fmt.Sprintf("var%d-1", k),
			test.var1,
		)
		var2, _ := vrb.NewVariable(
			fmt.Sprintf("var%d-2", k),
			test.var2,
		)
		actual, _ := vrb.Compare(var1, var2)
		expected := test.result
		if actual != expected {
			t.Errorf(
				"comparing %v and %v\nEXPECT: %v\n GET: %v\n",
				test.var1,
				test.var2,
				expected,
				actual,
			)
		}
	}
}
