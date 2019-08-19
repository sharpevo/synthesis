package instruction_test

import (
	"fmt"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"strings"
	"testing"
)

func TestInstructionVariableSet(t *testing.T) {
	testList := []struct {
		name    string
		value   string
		vrbtype vrb.VariableType
	}{
		{
			name:    "var1",
			value:   "string value",
			vrbtype: vrb.STRING,
		},
		{
			name:    "var2",
			value:   "2",
			vrbtype: vrb.INT,
		},
		{
			name:    "var3",
			value:   "3.0",
			vrbtype: vrb.FLOAT,
		},
		{
			name:    "var1",
			value:   "4",
			vrbtype: vrb.INT,
		},
	}
	i := instruction.InstructionVariableSet{}
	i.Env = interpreter.NewStack()
	_, err := i.Execute("var0")
	if err != nil {
		if !strings.Contains(err.Error(), "not enough") {
			t.Fatal(err)
		}
	}

	var var1Original *vrb.Variable

	for k, test := range testList {
		t.Logf("#%d", k)
		_, err := i.Execute(test.name, test.value)
		variable, found := i.Env.Get(test.name)
		if !found {
			t.Fatal(err)
		}
		if variable.Type != test.vrbtype {
			t.Errorf(
				"\nEXPECT: %v\nGET: %v\n",
				test.vrbtype,
				variable.Value,
			)
		}
		t.Logf("%#v", i.Env)
		if k == 0 {
			var1Original = variable
		}
	}

	variable, found := i.Env.Get("var1")
	if !found {
		t.Fatal(err)
	}
	expected := "4"
	actual := fmt.Sprintf("%d", variable.Value.(int64))
	if actual != expected {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			expected,
			actual,
		)
	}
	if variable.Type != vrb.INT {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			vrb.INT,
			variable.Type,
		)
	}
	if variable != var1Original {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			var1Original,
			variable,
		)
	}
}
