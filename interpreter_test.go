package interpreter_test

import (
	"fmt"
	"os"
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	interpreter.InitParser(TestInstructionMap)
	Test.SetTitle("PRINT")
	ret := m.Run()
	os.Exit(ret)
}

var TestInstructionMap = map[string]instruction.Instructioner{
	"TEST":   &Test,
	"PRINT":  &Test,
	"IMPORT": &instruction.Import,
	"ASYNC":  &instruction.Async,
	"RETRY":  &instruction.Retry,
	"MOVEX":  &instruction.MoveX,
}

type InstructionTest struct {
	instruction.Instruction
}

var Test InstructionTest

func (c *InstructionTest) Execute(args ...string) (interface{}, error) {
	return args[0] + "_test", nil
}

func TestParseLine(t *testing.T) {
	var tests = []struct {
		l string
		s *interpreter.Statement
	}{
		{
			l: `IMPORT C:\POSaM\scripts\async.script`,
			s: &interpreter.Statement{
				InstructionName: "IMPORT",
				Arguments:       []string{`C:\POSaM\scripts\async.script`},
			},
		},
		{
			l: `PRINT`,
			s: &interpreter.Statement{},
		},
		{
			l: `PRINT A B C`,
			s: &interpreter.Statement{
				InstructionName: "PRINT",
				Arguments:       []string{"A", "B", "C"},
			},
		},
	}

	for _, test := range tests {
		statement, _ := interpreter.ParseLine(test.l)
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
	suspend := false
	suspendTimer := time.NewTimer(2 * time.Second)
	recoverTimer := time.NewTimer(4 * time.Second)
	go func() {
		<-suspendTimer.C
		suspend = true
		<-recoverTimer.C
		suspend = false
	}()

	for _, test := range tests {
		statement, _ := interpreter.ParseLine(test.l)
		statementGroup := interpreter.StatementGroup{
			Execution: interpreter.SYNC,
			Stack:     concurrentmap.NewConcurrentMap(),
		}
		statement.StatementGroup = &statementGroup
		completec := make(chan interface{})
		go func() {
			<-completec
		}()
		resp := <-statement.Execute(terminatec, completec)
		if suspend {
			for {
				if !suspend {
					break
				}
			}
		}
		completec <- true
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
			f: "/home/yang/go/src/posam/interpreter/testscripts/script1",
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
			f: "/home/yang/go/src/posam/interpreter/testscripts/script3",
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
IMPORT /home/yang/go/src/posam/interpreter/testscripts/script2
PRINT 12
PRINT 13
RETRY -1 3
ASYNC /home/yang/go/src/posam/interpreter/testscripts/script4
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
ASYNC /home/yang/go/src/posam/interpreter/testscripts/scriptWithAsync
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
	timer := time.NewTimer(6 * time.Second)
	go func() {
		<-timer.C
		//close(terminatec)
	}()

	suspend := false
	suspendTimer := time.NewTimer(2 * time.Second)
	recoverTimer := time.NewTimer(4 * time.Second)
	go func() {
		<-suspendTimer.C
		//suspend = &suspended
		suspend = true
		fmt.Println(">", suspend)
		<-recoverTimer.C
		//suspend = &notSuspended
		suspend = false
		fmt.Println(">", suspend)
	}()

	for i, test := range tests {
		fmt.Printf("#%d\n", i)
		var resultList []string
		completec := make(chan interface{})
		go func() {
			<-completec
		}()
		switch test.s {
		case "file":
			statementGroup, _ := interpreter.ParseFile(
				test.f,
				interpreter.SYNC)
			statementGroup.Stack = concurrentmap.NewConcurrentMap()
			//resultList, _ = statementGroup.Execute(terminatec, nil)
			for resp := range statementGroup.Execute(terminatec, completec) {
				time.Sleep(1 * time.Second)
				if resp.Error != nil {
					fmt.Println(resp.Error)
					if resp.IgnoreError {
						resp.Completec <- true
					}
				} else {
					if suspend {
						for {
							if !suspend {
								break
							}
							time.Sleep(1 * time.Second)
						}
					}
					resp.Completec <- true
				}
				resultList = append(resultList, fmt.Sprintf("%v", resp.Output))
			}
		case "string":
			statementGroup := interpreter.StatementGroup{
				Execution: interpreter.SYNC,
				Stack:     concurrentmap.NewConcurrentMap(),
			}
			reader := strings.NewReader(test.f)
			interpreter.ParseReader(reader, &statementGroup)
			//resultList, _ = statementGroup.Execute(terminatec, nil)
			for resp := range statementGroup.Execute(terminatec, completec) {
				if resp.Error != nil {
					fmt.Println(resp.Error)
					if resp.IgnoreError {
						resp.Completec <- true
					}
				} else {
					if suspend {
						for {
							if !suspend {
								break
							}
							time.Sleep(1 * time.Second)
						}
					}
					resp.Completec <- true
				}
				resultList = append(resultList, fmt.Sprintf("%v", resp.Output))
			}
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
