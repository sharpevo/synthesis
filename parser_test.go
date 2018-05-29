package commandparser_test

import (
	"fmt"
	"posam/commandparser"
	"reflect"
	"sort"
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
			f: "/home/yang/go/src/posam/commandparser/testscripts/script1",
			r: []string{
				"11_test",
				"12_test",
				"21_test",
				"22_test",
				"13_test",
				"14_test",
			},
		},
		{
			f: "/home/yang/go/src/posam/commandparser/testscripts/script3",
			r: []string{
				"31_test",
				"32_test",
				"41_test",
				"42_test",
				"43_test",
				"44_test",
				"21_test",
				"22_test",
				"45_test",
				"46_test",
				"47_test",
				"33_test",
				"34_test",
			},
		},
	}

	for i, test := range tests {
		statementGroup, err := commandparser.ParseFile(
			test.f,
			commandparser.SYNC)
		if err != nil {
			fmt.Println(err)
		}
		resultList, _ := statementGroup.Execute()
		switch i {
		case 0:
			if !reflect.DeepEqual(resultList, test.r) {
				t.Errorf(
					"%d# EXPECT: %v\nGET:%v\n",
					i,
					test.r,
					resultList)
			}
		case 1:

			expect := append([]string{}, test.r...)
			get := append([]string{}, resultList...)
			sort.Strings(test.r)
			sort.Strings(resultList)
			if !reflect.DeepEqual(resultList, test.r) {
				t.Errorf(
					"%d# EXPECT: %v\nGET:%v\n",
					i,
					expect,
					get)
			}
		}
	}
}
