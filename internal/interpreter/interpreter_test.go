package interpreter_test

import (
	"fmt"
	"math/big"
	"os"
	"synthesis/internal/instruction"
	"synthesis/internal/interpreter"
	"synthesis/internal/interpreter/vrb"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var TestInstructionMap = make(interpreter.InstructionMapt)

type InstructionTest struct {
	instruction.Instruction
}

var Test InstructionTest

func (c *InstructionTest) Execute(args ...string) (interface{}, error) {
	return args[0] + "_test", nil
}

type InstructionSetStack struct {
	instruction.Instruction
}

func (c *InstructionSetStack) Execute(args ...string) (interface{}, error) {
	name := args[0]
	value := args[1]

	v, found := c.Env.Get(name)
	if found {
		v.Value = value
	} else {
		v, _ := vrb.NewVariable(name, value)
		c.Env.Set(v)
	}
	//c.Env.Set(args[0], &args[1])
	return nil, nil
}

func TestMain(m *testing.M) {
	TestInstructionMap.Set("TEST", InstructionTest{})
	TestInstructionMap.Set("PRINT", InstructionTest{})
	TestInstructionMap.Set("IMPORT", instruction.InstructionImport{})
	TestInstructionMap.Set("ASYNC", instruction.InstructionAsync{})
	TestInstructionMap.Set("RETRY", instruction.InstructionRetry{})
	TestInstructionMap.Set("MOVEX", instruction.InstructionMoveX{})
	TestInstructionMap.Set("SET", InstructionSetStack{})
	TestInstructionMap.Set("ADD", instruction.InstructionAdditionFloat64{})

	interpreter.InitParser(TestInstructionMap)
	Test.SetTitle("PRINT")
	ret := m.Run()
	os.Exit(ret)
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
			Stack:     interpreter.NewStack(),
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
			statementGroup.Stack = interpreter.NewStack()
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
				Stack:     interpreter.NewStack(),
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

func TestStack(t *testing.T) {
	sg := interpreter.StatementGroup{
		Execution: interpreter.SYNC,
		Stack:     interpreter.NewStack(),
	}

	var1, _ := vrb.NewVariable("KeyRoot1", "12.34")
	var2, _ := vrb.NewVariable("KeyRoot2", "stringtext")
	var3, _ := vrb.NewVariable("KeyRoot3", "22.11")
	sg.Stack.Set(var1)
	sg.Stack.Set(var2)
	sg.Stack.Set(var3)

	testString := `PRINT 10
SET KeyRoot2 abc
ADD KeyRoot1 43.21
IMPORT /home/yang/go/src/posam/interpreter/testscripts/scriptstack1
PRINT 15`
	reader := strings.NewReader(testString)
	interpreter.ParseReader(reader, &sg)
	terminatec := make(chan interface{})
	completec := make(chan interface{})
	var resultList []string
	go func() {
		<-completec
	}()
	for resp := range sg.Execute(terminatec, completec) {
		resultList = append(resultList, fmt.Sprintf("%v", resp.Output))
		resp.Completec <- true
	}

	if p, ok := sg.Stack.Get(var1.Name); ok {
		if p.Value.(*big.Float).String() != "55.55" {
			t.Errorf(
				"\nEXPECT: %v\nGET:%v\n",
				"55.55",
				p.Value,
			)
		}
	}
	if p, ok := sg.Stack.Get(var2.Name); ok {
		if p.Value != "fromimport" {
			t.Errorf(
				"\nEXPECT:%v\nGET:%v\n",
				"fromimport",
				p.Value,
			)
		}
	}
	if p, ok := sg.Stack.Get(var3.Name); ok {
		if p.Value.(*big.Float).String() != "33.33" {
			t.Errorf(
				"\nEXPECT: %v\nGET:%v\n",
				"33.33",
				p.Value,
			)
		}
	}
}
