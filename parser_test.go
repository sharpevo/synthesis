package commandparser_test

import (
	"fmt"
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

func TestParseFile(t *testing.T) {
	var tests = []struct {
		f string
		r []string
	}{
		{
			f: "/home/yang/go/src/posam/commandparser/script1",
			r: []string{
				"11_test",
				"12_test",
				"21_test",
				"22_test",
				"13_test",
				"14_test",
			},
		},
	}

	for _, test := range tests {
		statementGroup, err := commandparser.ParseFile(
			test.f,
			commandparser.SYNC)
		if err != nil {
			fmt.Println(err)
		}
		resultList, _ := statementGroup.Execute()
		if !reflect.DeepEqual(resultList, test.r) {
			t.Errorf(
				"EXPECT: %v\nGET:%v\n",
				test.r,
				resultList)
		}
	}
}
