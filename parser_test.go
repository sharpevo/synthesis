package commandparser_test

import (
	//"fmt"
	"os"
	"posam/commandparser"
	"reflect"
	"sort"
	"strings"
	"testing"
	//"time"
)

func TestMain(m *testing.M) {
	commandparser.InitParser(TestCommandMap)
	Test.SetTitle("PRINT")
	ret := m.Run()
	os.Exit(ret)
}

var TestCommandMap = map[string]commandparser.Commander{
	"TEST":   &Test,
	"PRINT":  &Test,
	"IMPORT": &commandparser.Import,
	"ASYNC":  &commandparser.Async,
	"RETRY":  &commandparser.Retry,
	"MOVEX":  &commandparser.MoveX,
}

type CommandTest struct {
	commandparser.Command
}

var Test CommandTest

func (c *CommandTest) Execute(args ...string) (interface{}, error) {
	return args[0] + "_test", nil
}

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

	terminatec := make(chan interface{})
	defer close(terminatec)
	for _, test := range tests {
		statement, _ := commandparser.ParseLine(test.l)
		resp := <-statement.Execute(terminatec)
		result := resp.Output
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
		s string
		f string
		r []string
	}{
		{
			s: "file",
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
			s: "file",
			f: "/home/yang/go/src/posam/commandparser/testscripts/script3",
			r: []string{
				"31_test",
				"32_test",
				"51_test",
				"52_test",
				"53_test",
				"54_test",
				"33_test",
				"34_test",
			},
		},
		{
			s: "string",
			f: `PRINT 10
PRINT 11
IMPORT /home/yang/go/src/posam/commandparser/testscripts/script2
PRINT 12
PRINT 13
RETRY -1 3
ASYNC /home/yang/go/src/posam/commandparser/testscripts/script4
PRINT 14
PRINT 15`,
			r: []string{
				"10_test",
				"11_test",
				"21_test",
				"22_test",
				"12_test",
				"13_test",
				"13_test",
				"13_test",
				"13_test",
				"41_test",
				"42_test",
				"43_test",
				"44_test",
				"45_test",
				"46_test",
				"47_test",
				"51_test",
				"52_test",
				"53_test",
				"54_test",
				"14_test",
				"15_test",
			},
		},
		{
			s: "string",
			f: `PRINT 10
PRINT 11
ASYNC /home/yang/go/src/posam/commandparser/testscripts/scriptWithAsync
PRINT 12
PRINT 13`,
			r: []string{
				"10_test",
				"11_test",
				"31_test",
				"32_test",
				"33_test",
				"34_test",
				"51_test",
				"52_test",
				"53_test",
				"54_test",
				"12_test",
				"13_test",
			},
		},
		{
			s: "string",
			f: `PRINT 10
MOVEX 5`,
			r: []string{
				"10_test",
				"Movable: 5",
			},
		},
	}

	terminatec := make(chan interface{})
	defer close(terminatec)

	//timer := time.NewTimer(10 * time.Nanosecond)
	//go func() {
	//<-timer.C
	//close(terminatec)
	//}()

	for i, test := range tests {
		var resultList []string
		switch test.s {
		case "file":
			statementGroup, _ := commandparser.ParseFile(
				test.f,
				commandparser.SYNC)
			resultList, _ = statementGroup.Execute(terminatec, nil)
		case "string":
			statementGroup := commandparser.StatementGroup{Execution: commandparser.SYNC}
			reader := strings.NewReader(test.f)
			commandparser.ParseReader(reader, &statementGroup)
			resultList, _ = statementGroup.Execute(terminatec, nil)
		}
		switch i {
		case 0:
			if !reflect.DeepEqual(resultList, test.r) {
				t.Errorf(
					"%d# EXPECT: %v\nGET:%v\n",
					i,
					test.r,
					resultList)
			}
		default:
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
