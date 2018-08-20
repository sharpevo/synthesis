package instruction_test

import (
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util/concurrentmap"
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
			value:   "new value",
			vrbtype: vrb.STRING,
		},
	}
	i := instruction.InstructionVariableSet{}
	i.Env = concurrentmap.NewConcurrentMap()
	_, err := i.Execute("var0")
	if err != nil {
		if !strings.Contains(err.Error(), "not enough") {
			t.Fatal(err)
		}
	}

	for k, test := range testList {
		t.Logf("#%d", k)
		_, err := i.Execute(test.name, test.value)
		v, found := i.Env.Get(test.name)
		if !found {
			t.Fatal(err)
		}
		variable, ok := v.(*vrb.Variable)
		if !ok {
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
	}

	v, found := i.Env.Get("var1")
	if !found {
		t.Fatal(err)
	}
	variable, ok := v.(*vrb.Variable)
	if !ok {
		t.Fatal(err)
	}
	expected := "new value"
	actual := variable.Value.(string)
	if actual != expected {
		t.Errorf(
			"\nEXPECT: %v\nGET: %v\n",
			expected,
			actual,
		)
	}
}
