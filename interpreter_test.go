package interpreter_test

import (
	"fmt"
	"os"
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"reflect"
	"sort"
	"strconv"
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
		v.(*interpreter.Variable).Value = value
	} else {
		v := interpreter.Variable{
			Name:  name,
			Value: value,
			Type:  "blah...",
		}
		c.Env.Set(name, &v)
	}
	//c.Env.Set(args[0], &args[1])
	return nil, nil
}

type InstructionAdd struct {
	instruction.Instruction
}

func (i *InstructionAdd) Execute(args ...string) (interface{}, error) {

	v1 := GetFloat64(args[0], i.Env)
	v2 := GetFloat64(args[1], i.Env)
	sum := v1 + v2
	v, _ := i.Env.Get(args[0])
	v.(*interpreter.Variable).Value = sum
	//*v.(*float64) = sum
	//i.Env.Set(args[0], v)
	return sum, nil
}

func GetFloat64(v interface{}, env *interpreter.Stack) float64 {
	switch i := v.(type) {
	case string:
		fmt.Println("parse string", i)
		v, found := env.Get(i)
		if !found {
			v, err := strconv.ParseFloat(i, 64)
			if err != nil {
				return float64(0)
			} else {
				return v
			}
		}
		return GetFloat64(v, env)
	case *interpreter.Variable:
		fmt.Println("parse variable", i)
		return GetFloat64(i.Value, env)
	case float64:
		fmt.Println("parse float", i)
		return i
	default:
		fmt.Println("parse none", i)
		return float64(0)
	}
}

func xGetFloat64(v interface{}, env *concurrentmap.ConcurrentMap) float64 {
	switch i := v.(type) {
	case float64:
		return float64(i)
	case string:
		vi, found := env.Get(string(i))
		if !found {
			f, err := strconv.ParseFloat(i, 64)
			if err != nil {
				return float64(0)
			}
			return f
		} else {
			//p := vi.(*float64)
			//return float64(*p)
			return GetFloat64(vi, env)
		}
	case *string:
		p, err := strconv.ParseFloat(*i, 64)
		if err != nil {
			return float64(0)
		} else {
			return p
		}
	case *float64:
		return float64(*i)
	default:
		return float64(0)
	}
}

func TestMain(m *testing.M) {
	TestInstructionMap.Set("TEST", InstructionTest{})
	TestInstructionMap.Set("PRINT", InstructionTest{})
	TestInstructionMap.Set("IMPORT", instruction.InstructionImport{})
	TestInstructionMap.Set("ASYNC", instruction.InstructionAsync{})
	TestInstructionMap.Set("RETRY", instruction.InstructionRetry{})
	TestInstructionMap.Set("MOVEX", instruction.InstructionMoveX{})
	TestInstructionMap.Set("SET", InstructionSetStack{})
	TestInstructionMap.Set("ADD", InstructionAdd{})

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
		Stack:     concurrentmap.NewConcurrentMap(),
	}

	v1 := interpreter.Variable{
		Name:  "KeyRoot1",
		Value: 12.34,
		Type:  "float64",
	}
	v2 := interpreter.Variable{
		Name:  "KeyRoot2",
		Value: "stringtext",
		Type:  "string",
	}
	v3 := interpreter.Variable{
		Name:  "KeyRoot3",
		Value: 22.11,
		Type:  "float64",
	}
	sg.Stack.Set("KeyRoot1", &v1)
	sg.Stack.Set("KeyRoot2", &v2)
	sg.Stack.Set("KeyRoot3", &v3)

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

	if v, ok := sg.Stack.Get("KeyRoot1"); ok {
		p := v.(*interpreter.Variable)
		if p.Value != 55.55 {
			t.Errorf(
				"\nEXPECT: %v\nGET:%v\n",
				"12.34",
				p.Value,
			)
		}
	}
	if v, ok := sg.Stack.Get("KeyRoot2"); ok {
		p := v.(*interpreter.Variable)
		if p.Value != "fromimport" {
			t.Errorf(
				"\nEXPECT:%v\nGET:%v\n",
				"fromimport",
				p.Value,
			)
		}
	}
	if v, ok := sg.Stack.Get("KeyRoot3"); ok {
		p := v.(*interpreter.Variable)
		if p.Value != 33.33 {
			t.Errorf(
				"\nEXPECT: %v\nGET:%v\n",
				"33.33",
				p.Value,
			)
		}
	}
}
