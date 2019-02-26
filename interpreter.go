package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"posam/config"
	"posam/dao"
	"posam/gui/uiutil"
	"posam/instruction"
	"posam/interpreter/vrb"
	"posam/util"
	"posam/util/blockingqueue"
	"runtime"
	"strings"
	"sync"
	"time"
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

type Statement struct {
	//sync.RWMutex
	Stack *instruction.Stack
	//Title string
	InstructionName string
	Arguments       []string
	IgnoreError     bool
}

type StatementGroup struct {
	sync.Mutex
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

func InitParser(instructionMap *dao.InstructionMapt) error {
	return nil
}

func RunInstruction(
	name string,
	arguments []string,
	//stack *interpreter.Stack,
	stack *instruction.Stack,
) (resp interface{}, err error) {

	fmt.Println("### running instruction", name, arguments)
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
	case "GETFLOATBYINDEX":
		i := instruction.InstructionVariableGetFloatByIndex{}
		//i.Env = instruction.NewStack(stack)
		i.Env = stack
		resp, err = i.Execute(arguments...)
	case "GETINTBYINDEX":
		i := instruction.InstructionVariableGetIntByIndex{}
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
		s.InstructionName, s.Arguments, s.Stack)
	resp.Output = output
	resp.Error = err
	resp.Completec = completec
	resp.IgnoreError = s.IgnoreError
	fmt.Println("instructionname", s.InstructionName)
	fmt.Println("arguments", s.Arguments)
	fmt.Println("output", output)
	log.Printf("'%s: %s' produces %q\n", s.InstructionName, s.Arguments, output)
	message := ""
	if resp.Error != nil {
		message = resp.Error.Error()
	}
	fmt.Println(message)
	if GCABLE {
		runtime.GC()
	}
	return resp
}

func (s *Statement) Execute(terminatec <-chan interface{}, completec chan<- interface{}) <-chan Response {

	respc := make(chan Response)
	util.Go(func() {
		defer close(respc)
		fmt.Println("### executing instruction", s.InstructionName, s.Arguments)
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
	runtime.KeepAlive(outputc)
	return outputc
}

func bridge(terminatec <-chan interface{}, chanc <-chan <-chan Response) <-chan Response {
	valStream := make(chan Response)
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
	respcc := make(chan (<-chan Response))

	util.Go(func() {
		var wg sync.WaitGroup
		wg.Add(g.ItemList.Length())
		defer close(respcc)

		for i := 0; i < g.ItemList.Length(); i++ {
			completec := make(chan interface{})
			go func(index int) {
				g.ItemList.Lock()
				itemi, err := g.ItemList.GetLockless(index)
				script, _ := itemi.(*Script)
				itemInterface, err := script.ParseUnit(g.Stack)
				if err != nil {
					uiutil.App.ShowMessageSlot(err.Error())
				}
				g.ItemList.Unlock()
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
					log.Printf("NO MATCH %T: %#v!\n", t, itemInterface)
				}
			}(i)

			go func() {
				defer wg.Done()
				<-completec
			}()
		}

		wg.Wait()
		pcompletec <- true
	})

	return respcc
}

func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}, pcompletec chan<- interface{}) <-chan <-chan Response { // {{{

	log.Println("==== SYNC ====")
	respcc := make(chan (<-chan Response))

	util.Go(func() {
		defer close(respcc)
		list := g.ItemList.ItemList()
		length := g.ItemList.Length()
		for i := 0; i < length; i++ {
			completec := make(chan interface{})
			g.ItemList.Lock()
			index := i
			itemi := list[index]
			script, _ := itemi.(*Script)
			itemInterface, err := script.ParseUnit(g.Stack)
			if err != nil {
				uiutil.App.ShowMessageSlot(err.Error())
			}
			g.ItemList.Unlock()
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				respcc <- item.Execute(terminatec, completec)
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				respcc <- item.Execute(terminatec, completec)
			default:
				log.Printf("!!!!!! NO MATCH %T: %#v!\n", t, itemInterface)
			}
			<-completec
			stack := g.Stack
			cmap, _ := stack.Get("SYS_CUR")
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
				varCur.SetValue(int64(0))
				cmap.Unlock()
				continue
			} else if cur < 0 {
				varCur.SetValue(int64(0))
				cmap.Unlock()
				break
			}
			cmap.Unlock()
		}
		pcompletec <- true
		if GCABLE {
			runtime.GC()
		}
	})

	return respcc
} // }}}

func (g *StatementGroup) Execute(terminatec <-chan interface{}, completec chan<- interface{}) <-chan Response {
	switch g.Execution {
	case SYNC:
		return bridge(terminatec, g.ExecuteSync(terminatec, completec))
	case ASYNC:
		return bridge(terminatec, g.ExecuteAsync(terminatec, completec))
	default:
		msg := fmt.Sprintf("!!!!!! Execution type is not matched\n%#v\nSYNC:%v\nASYNC:%v\n", g, SYNC, ASYNC)
		fmt.Println(msg)
		uiutil.App.ShowMessageSlot(msg)
		resultc := make(chan Response)
		go func() {
			defer close(resultc)
			<-time.After(time.Second)
			resp := Response{}
			resp.Completec = completec
			resultc <- resp
		}()
		return resultc
	}
}

type Script struct {
	sync.Mutex
	Line    string
	Scripts []*Script
}

var RootScript *Script

func (s *Script) ParseReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		script := &Script{}
		script.Line = line
		script.ParseItem()
		s.Scripts = append(s.Scripts, script)
	}
	err := scanner.Err()
	if err != nil {
		fmt.Println("!!! Error:", err)
	}
	return err
}

func (s *Script) ParseItem() error {
	s.Lock()
	defer s.Unlock()
	list := strings.Split(s.Line, " ")
	if len(list) < 2 {
		return fmt.Errorf("Error: %s", "Invalid syntax")
	}
	instruction := list[0]
	arguments := list[1:]
	switch instruction {
	case "IMPORT":
		return s.ParseFile(arguments[0])
	case "ASYNC":
		return s.ParseFile(arguments[0])
	}
	return fmt.Errorf("invalid instruction")
}

func (s *Script) ParseFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	return s.ParseReader(file)
}

func (s *Script) String() string {
	result := fmt.Sprintf("%q: ", s.Line)
	for _, script := range s.Scripts {
		result = fmt.Sprintf("%v\n%v", result, script.String())
	}
	return result
}

func (s *Script) ParseUnit(stack *instruction.Stack) (interface{}, error) {
	fmt.Println("ParseUnit")
	item, _ := s.ParseLine()
	s.Lock()
	defer s.Unlock()
	fmt.Println(">> parse unit", s)
	switch item.InstructionName {
	case "IMPORT":
		return s.CreateStatementGroup(stack, SYNC)
	case "ASYNC":
		return s.CreateStatementGroup(stack, ASYNC)
	default:
		item.Stack = stack
		return item, nil
	}
}

func (s *Script) ParseLine() (*Statement, error) {
	s.Lock()
	defer s.Unlock()
	itemList := strings.Split(s.Line, " ")
	statement := &Statement{}
	if len(itemList) < 2 {
		return statement, fmt.Errorf("Error: %s", "Invalid syntax")
	}
	statement.InstructionName = itemList[0]
	statement.Arguments = itemList[1:]
	return statement, nil
}

func (s *Script) CreateStatementGroup(
	stack *instruction.Stack,
	execution ExecutionType,
) (*StatementGroup, error) {
	statementGroup := &StatementGroup{}
	statementGroup.Execution = execution
	statementGroup.ItemList = blockingqueue.NewBlockingQueue()
	for _, script := range s.Scripts {
		script.Lock()
		statementGroup.ItemList.Append(script) // create at runtime
		script.Unlock()
	}
	fmt.Printf("create statementgroup: %#v\n", statementGroup.ItemList.ItemList())
	statementGroup.Stack = instruction.NewStack(stack)
	return statementGroup, nil
}
