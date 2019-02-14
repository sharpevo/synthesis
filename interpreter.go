package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"posam/config"
	"posam/dao"
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"reflect"
	//"runtime"
	"strconv"
	"strings"
	"sync"
)

func init() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	log.SetFlags(log.LstdFlags)
}

type ExecutionType int

const (
	SYNC ExecutionType = iota
	ASYNC
)

var (
	GCABLE = config.GetBool("general.gc.manual")
)

//type InstructionMapt map[string]reflect.Type

//var InstructionMap = make(InstructionMapt)

//func (m InstructionMapt) Set(k string, v interface{}) {
//m[k] = reflect.TypeOf(v)
//}

//func (m InstructionMapt) Get(name string) (reflect.Value, error) {
//if t, ok := m[name]; ok {
//return reflect.New(t), nil
//}
//return reflect.Value{}, fmt.Errorf("instruction %q not found", name)
//}

type InstMapt struct {
	cmap *concurrentmap.ConcurrentMap
}

var InstMap = InstMapt{
	cmap: concurrentmap.NewConcurrentMap(),
}

func (i *InstMapt) AddType(k string, v interface{}) {
	value, ok := v.(reflect.Type)
	if ok {
		i.cmap.Set(k, value)
	} else {
		i.cmap.Set(k, reflect.TypeOf(v))
	}
}
func (i *InstMapt) Lock() {
	i.cmap.Lock()
}
func (i *InstMapt) Unlock() {
	i.cmap.Unlock()
}
func (i *InstMapt) GetLockless(key string) (interface{}, bool) {
	return i.cmap.GetLockless(key)
}

type InstructionMapt struct {
	lock sync.RWMutex
	cmap *concurrentmap.ConcurrentMap
}

func NewInstructionMap() *InstructionMapt {
	return &InstructionMapt{
		cmap: concurrentmap.NewConcurrentMap(),
	}
}

var InstructionMap = NewInstructionMap()

func (m *InstructionMapt) Set(k string, v interface{}) {
	//value := v.(reflect.Type)
	value, ok := v.(reflect.Type)
	if ok {
		m.cmap.Set(k, value)
	} else {
		m.cmap.Set(k, reflect.TypeOf(v))
	}
}

func (m *InstructionMapt) Lock() {
	m.cmap.Lock()
}

func (m *InstructionMapt) Unlock() {
	m.cmap.Unlock()
}

//func (m *InstructionMapt) Get(name string) (reflect.Value, error) {
//if tv, ok := m.cmap.Get(name); ok {
//t := tv.(reflect.Type)
//v := reflect.New(t)
//fmt.Printf("get %#v, %v\n", v, ok)
//return v, nil
//}
//return reflect.Value{}, fmt.Errorf("instruction %q not found", name)
//}

func (m *InstructionMapt) GetNew(name string) (reflect.Type, error) {
	if tv, ok := m.cmap.Get(name); ok {
		t := tv.(reflect.Type)
		return t, nil
	}
	return nil, fmt.Errorf("instruction %q not found", name)
}

func (m *InstructionMapt) Get(name string) (*InstructionInterface, error) {
	if tv, ok := m.cmap.Get(name); ok {
		t := tv.(reflect.Type)
		return &InstructionInterface{
			value: reflect.New(t),
		}, nil
	}
	return nil, fmt.Errorf("instruction %q not found", name)
}

func (m *InstructionMapt) GetLockless(name string) (reflect.Type, error) {
	if tv, ok := m.cmap.GetLockless(name); ok {
		t := tv.(reflect.Type)
		return t, nil
	}
	return nil, fmt.Errorf("instruction %q not found", name)
}

func (m *InstructionMapt) GetLocklessDepre(name string) (*InstructionInterface, error) {
	if tv, ok := m.cmap.GetLockless(name); ok {
		t := tv.(reflect.Type)
		i := &InstructionInterface{}
		i.SetValue(reflect.New(t))
		return i, nil
	}
	return nil, fmt.Errorf("instruction %q not found", name)
}

func (m *InstructionMapt) Iter() <-chan concurrentmap.Item {
	return m.cmap.Iter()
}

type InstructionInterface struct {
	lock  sync.RWMutex
	value reflect.Value
}

func (i *InstructionInterface) Value() reflect.Value {
	//i.lock.Lock()
	//i.lock.Unlock()
	return i.value
}

func (i *InstructionInterface) GetValueLockless() reflect.Value {
	return i.value
}

func (i *InstructionInterface) SetValue(value reflect.Value) {
	//i.lock.Lock()
	//defer i.lock.Unlock()
	i.value = value
}

func (i *InstructionInterface) Lock() {
	i.lock.Lock()
}

func (i *InstructionInterface) Unlock() {
	i.lock.Unlock()
}

type Statement struct {
	sync.RWMutex
	StatementGroup  *StatementGroup
	InstructionName string
	Arguments       []string
	IgnoreError     bool
}

type StatementGroup struct {
	Execution ExecutionType
	ItemList  *blockingqueue.BlockingQueue
	Stack     *instruction.Stack
}

type Response struct {
	Error       error
	Output      interface{}
	Completec   chan<- interface{}
	IgnoreError bool
}

type Info struct {
	Line   int
	Column int
}

//func InitParser(instructionMap InstructionMapt) error {
//for k, v := range instructionMap {
//InstructionMap[k] = v
//}
//return nil
//}

func InitParser(instructionMap *dao.InstructionMapt) error {
	//InstructionMap = instructionMap
	for item := range instructionMap.Iter() {
		k := item.Key
		v := item.Value.(reflect.Type)
		InstructionMap.Set(k, v)
		InstMap.AddType(k, v)
	}
	return nil
}

func ParseLine(line string) (*Statement, error) {
	itemList := strings.Split(line, " ")
	statement := &Statement{}
	if len(itemList) < 2 {
		return statement, fmt.Errorf("Error: %s", "Invalid syntax")
	}
	statement.InstructionName = itemList[0]
	statement.Arguments = itemList[1:]
	return statement, nil
}

func RunInstruction(
	name string,
	arguments []string,
	//stack *interpreter.Stack,
	stack *instruction.Stack,
) (resp interface{}, err error) {
	switch name {

	// UNKNOWN{{{
	case "ADD":
		i := instruction.InstructionAddition{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "SUB":
		i := instruction.InstructionSubtraction{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "DIV":
		i := instruction.InstructionDivision{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MUL":
		i := instruction.InstructionMultiplication{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOD":
		i := instruction.InstructionModulo{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "EQGOTO":
		i := instruction.InstructionControlFlowEqualGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "ERRGOTO":
		i := instruction.InstructionControlFlowErrGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "GOTO":
		i := instruction.InstructionControlFlowGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "GTGOTO":
		i := instruction.InstructionControlFlowGreaterThanGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "LOOP":
		i := instruction.InstructionControlFlowLoop{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "LTGOTO":
		i := instruction.InstructionControlFlowLessThanGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "NEGOTO":
		i := instruction.InstructionControlFlowNotEqualGoto{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "RETURN":
		i := instruction.InstructionControlFlowReturn{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "PRINT":
		i := instruction.InstructionPrint{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "SLEEP":
		i := instruction.InstructionSleep{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "CMPVAR":
		i := instruction.InstructionVariableCompare{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "GETVAR":
		i := instruction.InstructionVariableGet{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "GETVARBYINDEX":
		i := instruction.InstructionVariableGetByIndex{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "SETVAR":
		i := instruction.InstructionVariableSet{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
		// }}}

		// CAN{{{
	case "CANMOVEABS":
		i := instruction.InstructionCANMotorMoveAbsolute{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "CANMOVEREL":
		i := instruction.InstructionCANMotorMoveRelative{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "CANRESETMOTOR":
		i := instruction.InstructionCANMotorReset{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "HUMITURE":
		i := instruction.InstructionCANSensorHumiture{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "OXYGENCONC":
		i := instruction.InstructionCANSensorOxygenConc{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "PRESSURE":
		i := instruction.InstructionCANSensorPressure{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "SWITCH":
		i := instruction.InstructionCANSwitcherControl{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "SWITCHCOND":
		i := instruction.InstructionCANSwitcherControlAdvanced{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "ROMREAD":
		i := instruction.InstructionCANSystemRomRead{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "ROMWRITE":
		i := instruction.InstructionCANSystemRomWrite{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
		// }}}

		// PRINTERHEAD{{{
	case "ERRORCODE":
		i := instruction.InstructionPrinterHeadErrorCode{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "PRINTDATA":
		i := instruction.InstructionPrinterHeadPrintData{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "PRINTERSTATUS":
		i := instruction.InstructionPrinterHeadPrinterStatus{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "WAVEFORM":
		i := instruction.InstructionPrinterHeadWaveform{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	// }}}

	// RINTER{{{
	case "LOADCYCLE":
		i := instruction.InstructionPrinterLoadCycle{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "LOADEXEC":
		i := instruction.InstructionPrinterLoadExec{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "LOADGROUP":
		i := instruction.InstructionPrinterLoadFormation{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
		// }}}

		// TML{{{
	case "MOVEABS":
		i := instruction.InstructionTMLMoveAbs{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOVEABSX":
		i := instruction.InstructionTMLMoveAbsX{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOVEABSY":
		i := instruction.InstructionTMLMoveAbsY{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOVEREL":
		i := instruction.InstructionTMLMoveRel{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOVERELX":
		i := instruction.InstructionTMLMoveRelX{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "MOVERELY":
		i := instruction.InstructionTMLMoveRelY{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "POSITIONX":
		i := instruction.InstructionTMLPositionX{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "POSITIONY":
		i := instruction.InstructionTMLPositionY{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
		// }}}

	default:
		err = fmt.Errorf("invalid instruction %s", name)
	}
	return
}

func (s *Statement) Run(completec chan<- interface{}) Response {
	resp := Response{}
	output, err := RunInstruction(
		s.InstructionName, s.Arguments, s.StatementGroup.Stack)
	resp.Output = output
	resp.Error = err
	resp.Completec = completec
	resp.IgnoreError = s.IgnoreError
	log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
	message := ""
	if resp.Error != nil {
		message = resp.Error.Error()
	}
	fmt.Println(message)
	return resp
}

func (s *Statement) RunNewDepre(completec chan<- interface{}) Response {
	var err error
	resp := Response{}
	done := make(chan struct{})
	InstMap.Lock()
	instTypei, found := InstMap.GetLockless(s.InstructionName)
	if !found {
		resp.Error = fmt.Errorf("invalid instruction %s", s.InstructionName)
		InstMap.Unlock()
		return resp
	}
	instType, _ := instTypei.(reflect.Type)
	instValue := reflect.New(instType)
	fmt.Printf("ask for value: %#v\n", instValue)
	reflect.Indirect(instValue).FieldByName("Env").Set(reflect.ValueOf(s.StatementGroup.Stack))
	method := instValue.MethodByName("Execute")
	arguments := make([]reflect.Value, len(s.Arguments))
	for i, _ := range s.Arguments {
		arguments[i] = reflect.ValueOf(s.Arguments[i])
	}
	var outputValueList []reflect.Value
	util.Go(func() {
		outputValueList = method.Call(arguments)
		output := outputValueList[0].Interface()
		erri := outputValueList[1].Interface()
		if erri != nil {
			err = erri.(error)
		} else {
			err = nil
		}

		resp.Output = output
		resp.Error = err
		resp.Completec = completec
		resp.IgnoreError = s.IgnoreError
		log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
		message := ""
		if resp.Error != nil {
			message = resp.Error.Error()
		}
		fmt.Println(message)
		done <- struct{}{}
	})
	InstMap.Unlock()
	<-done
	return resp
}

//func (s *Statement) Run(completec chan<- interface{}) Response {
//resp := Response{}
//done := make(chan struct{})
////InstructionMap.Lock()
//InstructionMap.lock.Lock()
//InstructionMap.Lock()
//instructionInstancev, err := InstructionMap.GetLocklessDepre(s.InstructionName)
//instructionInstancev.Lock()
////instructionInstancev, err := InstructionMap.Get(s.InstructionName)
//if err != nil {
//resp.Error = err
////InstructionMap.Unlock()
//InstructionMap.lock.Unlock()
//return resp
//}
//value := instructionInstancev.Value()
////s.Lock() // lock the execution
////value := reflect.New(instructionInstancev)
//fmt.Printf("ask for value: %#v\n", value)
//method := value.MethodByName("Execute")
//reflect.Indirect(value).FieldByName("Env").Set(reflect.ValueOf(s.StatementGroup.Stack))
//arguments := make([]reflect.Value, len(s.Arguments))
//for i, _ := range s.Arguments {
//arguments[i] = reflect.ValueOf(s.Arguments[i])
//}
//var outputValueList []reflect.Value
//go func() {
//outputValueList = method.Call(arguments)
////InstructionMap.lock.Unlock() // not well
////s.Unlock()
//output := outputValueList[0].Interface()
//erri := outputValueList[1].Interface()
//if erri != nil {
//err = erri.(error)
//} else {
//err = nil
//}

//resp.Output = output
//resp.Error = err
//resp.Completec = completec
//resp.IgnoreError = s.IgnoreError
//log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
//message := ""
//if resp.Error != nil {
//message = resp.Error.Error()
//}
//fmt.Println(message)
//done <- struct{}{}
//}()
//instructionInstancev.Unlock()
//InstructionMap.Unlock()
//InstructionMap.lock.Unlock()
//<-done
//return resp
//}

// {{{
func (s *Statement) RunDepreciated2(completec chan<- interface{}) Response {
	resp := Response{}
	done := make(chan struct{})
	found := false
	var err error
	for item := range InstructionMap.Iter() {
		if found {
			continue
		}
		if item.Key == s.InstructionName {
			instructionValue := &InstructionInterface{
				value: reflect.New(item.Value.(reflect.Type)), //.(*InstructionInterface)
			}
			reflect.Indirect(instructionValue.Value()).FieldByName("Env").Set(reflect.ValueOf(s.StatementGroup.Stack))
			arguments := make([]reflect.Value, len(s.Arguments))
			for i, _ := range s.Arguments {
				arguments[i] = reflect.ValueOf(s.Arguments[i])
			}
			method := instructionValue.Value().MethodByName("Execute")
			var outputValueList []reflect.Value
			util.Go(func() {
				outputValueList = method.Call(arguments)
				output := outputValueList[0].Interface()
				erri := outputValueList[1].Interface()
				if erri != nil {
					err = erri.(error)
				} else {
					err = nil
				}

				resp.Output = output
				resp.Error = err
				resp.Completec = completec
				resp.IgnoreError = s.IgnoreError
				log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
				message := ""
				if resp.Error != nil {
					message = resp.Error.Error()
				}
				//instructionInstancev.MethodByName("IssueError").Call([]reflect.Value{reflect.ValueOf(message)})
				fmt.Println(message)
				done <- struct{}{}
			})
		}
	}
	<-done
	return resp
}

func (s *Statement) RunDepreciated(completec chan<- interface{}) Response {
	resp := Response{}
	//if _, ok := InstructionMap[s.InstructionName]; !ok {
	if _, err := InstructionMap.Get(s.InstructionName); err != nil {
		resp.Error = fmt.Errorf("Invalid instruction %q", s.InstructionName)
	} else {
		instructionInstancev, err := InstructionMap.Get(s.InstructionName)
		if err != nil {
			resp.Error = err
			return resp
		}
		fmt.Println("#### instructionInstancev", instructionInstancev)
		reflect.Indirect(instructionInstancev.Value()).FieldByName("Env").Set(reflect.ValueOf(s.StatementGroup.Stack))

		arguments := make([]reflect.Value, len(s.Arguments))
		for i, _ := range s.Arguments {
			arguments[i] = reflect.ValueOf(s.Arguments[i])
		}

		//outputValueList := instructionInstancev.MethodByName("Execute").Call(arguments)
		fmt.Println("#### check instructionInstancev", instructionInstancev)
		method := instructionInstancev.Value().MethodByName("Execute")
		fmt.Println("#### parse method", method, s.InstructionName, arguments)
		outputValueList := method.Call(arguments)

		output := outputValueList[0].Interface()
		erri := outputValueList[1].Interface()
		if erri != nil {
			err = erri.(error)
		} else {
			err = nil
		}

		resp.Output = output
		resp.Error = err
		resp.Completec = completec
		resp.IgnoreError = s.IgnoreError
		log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
		message := ""
		if resp.Error != nil {
			message = resp.Error.Error()
		}
		//instructionInstancev.MethodByName("IssueError").Call([]reflect.Value{reflect.ValueOf(message)})
		fmt.Println(message)
		// not in the current reflection but parent
		//method := instructionInstancev.MethodByName("IssueError")
		//if method.IsValid() {
		//method.Call([]reflect.Value{reflect.ValueOf(message)})
		//} else {
		//log.Println(">>>>>>>>>>>>>>>>>> invalid IssueError")
		//}
	}
	return resp
}

// }}}

func (s *Statement) Execute(terminatec <-chan interface{}, completec chan<- interface{}) <-chan Response {

	respc := make(chan Response)
	//respc := make(chan Response, 1)
	util.Go(func() {
		defer close(respc)
		//defer runtime.GC()
		for {
			select {
			case <-terminatec:
				resp := Response{}
				log.Printf("Termiante '%s: %s'\n\n", s.InstructionName, s.Arguments)
				resp.Error = fmt.Errorf("Terminated %q", s.InstructionName)
				respc <- resp
				return
			case respc <- s.Run(completec):
				log.Printf("'%s: %s' done\n", s.InstructionName, s.Arguments)
				return
			}
		}
	})
	log.Printf("'%s: %s' execute thread exits\n", s.InstructionName, s.Arguments)
	//time.Sleep(1 * time.Second)
	return respc
}

func tryRead(terminatec <-chan interface{}, inputc <-chan Response) <-chan Response {
	outputc := make(chan Response)
	util.Go(func() {
		defer close(outputc)
		for {
			select {
			case <-terminatec:
				return
			case resp, ok := <-inputc:
				if ok == false {
					return
				}
				select {
				case outputc <- resp:
				case <-terminatec:
				}
			}
		}
	})
	return outputc
}

func bridge(terminatec <-chan interface{}, chanc <-chan <-chan Response) <-chan Response {
	valStream := make(chan Response)
	//valStream := make(chan Response, 1)
	util.Go(func() {
		defer close(valStream)
		for {
			var stream <-chan Response
			select {
			case <-terminatec:
				return
			case mayStream, ok := <-chanc:
				if ok == false {
					return
				}
				stream = mayStream
			}
			util.Go(func() {
				for val := range tryRead(terminatec, stream) {
					select {
					case valStream <- val:
					case <-terminatec:
					}
				}
			})

		}
	})
	return valStream
}

func (g *StatementGroup) ExecuteAsync(terminatec <-chan interface{}, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== ASYNC ====")
	//respcc := make(chan (<-chan Response), len(g.ItemList))
	respcc := make(chan (<-chan Response))

	util.Go(func() {
		var wg sync.WaitGroup
		wg.Add(g.ItemList.Length())
		defer close(respcc)

		for itemi := range g.ItemList.Iter() {
			itemInterface := itemi.Value
			completec := make(chan interface{})
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				util.Go(func() {
					respcc <- item.Execute(terminatec, completec)
				})
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				util.Go(func() {
					respcc <- item.Execute(terminatec, completec)
				})
			default:
				log.Printf("NO MATCH %T!\n", t)
			}

			go func(itemInterface interface{}) {
				defer wg.Done()
				<-completec
				fmt.Println("complete", itemInterface)
			}(itemInterface)
		}

		wg.Wait()
		pcompletec <- true
	})

	return respcc
}

func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== SYNC ====")
	respcc := make(chan (<-chan Response))

	util.Go(func() {
		defer close(respcc)

		length := g.ItemList.Length()
		itemList := g.ItemList
		for i := 0; i < length; i++ {
			completec := make(chan interface{})
			//itemInterface, err := g.ItemList.Get(i)
			itemInterface, err := itemList.Get(i)
			if err != nil {
				log.Println(err)
				continue
			}
			switch t := itemInterface.(type) {
			case Statement, *Statement:

				item, _ := itemInterface.(*Statement)
				if item.InstructionName == "RETRY" &&
					// TODO: previous error
					true {
					lineNum, _ := strconv.Atoi(item.Arguments[0])
					count, _ := strconv.Atoi(item.Arguments[1])
					if count < 1 {
						// TODO: panic
						log.Printf(
							"Failed to retry at line %d in %d attempts\n",
							lineNum,
							count)
						continue
					}
					newLineIndex := i + lineNum
					if newLineIndex < 0 {
						i = 0
					} else {
						i = newLineIndex
					}
					count -= 1
					item.Arguments[1] = strconv.Itoa(count)
					i -= 1 //trade off the i++
					continue
				} else {
					if i < length-1 {
						nextItem, err := g.ItemList.Get(i + 1)
						if err != nil {
							log.Println(err)
						}
						if s, ok := nextItem.(*Statement); ok &&
							s.InstructionName == "RETRY" {
							item.IgnoreError = true
						}
					}
					respcc <- item.Execute(terminatec, completec)
					//log.Printf("'%s: %s' complet", item.InstructionName, item.Arguments)
				}

				if i < length-1 {
					nextItem, err := g.ItemList.Get(i + 1)
					if err != nil {
						log.Println(err)
					}
					if s, ok := nextItem.(*Statement); ok &&
						s.InstructionName == "ERRGOTO" {
						item.IgnoreError = true
					}
				}

			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				respcc <- item.Execute(terminatec, completec)
			default:
				log.Printf("NO MATCH %T!\n", t)
			}

			<-completec
			//varCur, _ := g.Stack.Get("SYS_CUR")
			//cur64 := varCur.GetValue().(int64)
			//cur := int(cur64)
			//if cur > 0 {
			//target := cur - 2
			//if target == i {
			//target += 1
			//}
			//i = target
			////varCur.Value = int64(0)
			//varCur.SetValue(int64(0))
			//} else if cur < 0 {
			////varCur.Value = int64(0)
			//varCur.SetValue(int64(0))
			//break
			//}
			cmap, _ := g.Stack.Get("SYS_CUR")
			cmap.Lock()
			curv, _ := cmap.GetLockless("SYS_CUR")
			varCur, _ := curv.(*vrb.Variable)
			cur64 := varCur.GetValue().(int64)
			cur := int(cur64)
			if cur > 0 {
				target := cur - 2
				if target == i {
					target += 1
				}
				i = target
				//varCur.Value = int64(0)
				varCur.SetValue(int64(0))
				//cmap.Unlock()
			} else if cur < 0 {
				//varCur.Value = int64(0)
				varCur.SetValue(int64(0))
				cmap.Unlock()
				break
			}
			cmap.Unlock()
		}

		pcompletec <- true
	})

	return respcc
}

func (g *StatementGroup) Execute(terminatec <-chan interface{}, completec chan<- interface{}) <-chan Response {
	switch g.Execution {
	case SYNC:
		return bridge(terminatec, g.ExecuteSync(terminatec, completec))
	case ASYNC:
		return bridge(terminatec, g.ExecuteAsync(terminatec, completec))
	}
	resultc := make(chan Response)
	close(resultc)
	return resultc
}

func ParseFile(
	//stack *Stack,
	stack *instruction.Stack,
	filePath string,
	execution ExecutionType) (*StatementGroup, error) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	statementGroup := StatementGroup{
		Execution: execution,
		ItemList:  blockingqueue.NewBlockingQueue(),
		Stack:     stack,
	}

	return ParseReader(file, &statementGroup)
}

func ParseReader(reader io.Reader, statementGroup *StatementGroup) (*StatementGroup, error) {
	if statementGroup.Stack == nil {
		statementGroup.Stack = instruction.NewStack()
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		statement, _ := ParseLine(line)
		statement.StatementGroup = statementGroup
		switch statement.InstructionName {
		case "IMPORT":
			subStatementGroup, _ := ParseFile(
				instruction.NewStack(statementGroup.Stack),
				statement.Arguments[0],
				SYNC)
			statementGroup.ItemList.Append(subStatementGroup)
		case "ASYNC":
			subStatementGroup, _ := ParseFile(
				instruction.NewStack(statementGroup.Stack),
				statement.Arguments[0],
				ASYNC)
			statementGroup.ItemList.Append(subStatementGroup)
		default:
			statementGroup.ItemList.Append(statement)
		}
	}
	return statementGroup, nil

}
