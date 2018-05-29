package commandparser_test

import (
	"posam/commandparser"
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	var tests = []struct {
		l string
		s *commandparser.Statement
	}{
		{
			l: `IMPORT C:\POSaM\scripts\async.script`,
			s: &commandparser.Statement{
				CommandName: "IMPORT",
				Arguments:   []string{`C:\POSaM\scripts\async.script`},
			},
		},
		{
			l: `PRINT`,
			s: &commandparser.Statement{},
		},
		{
			l: `PRINT A B C`,
			s: &commandparser.Statement{
				CommandName: "PRINT",
				Arguments:   []string{"A", "B", "C"},
			},
		},
	}

	for _, test := range tests {
		statement, _ := commandparser.ParseLine(test.l)
		if !reflect.DeepEqual(statement, test.s) {
			t.Errorf(
				"EXPECT: %v\n GET: %v\n\n",
				test.s,
				statement)
		}
	}
}

func TestExecute(t *testing.T) {
	var tests = []struct {
		l string
		r string
	}{
		{
			l: `TEST ABC DEF GHI`,
			r: `ABC_test`,
		},
	}

	for _, test := range tests {
		statement, _ := commandparser.ParseLine(test.l)
		result, _ := statement.Execute()
		if result != test.r {
			t.Errorf(
				"EXPECT: %v\n GET: %v\n\n",
				test.r,
				result)

		}
	}
}
