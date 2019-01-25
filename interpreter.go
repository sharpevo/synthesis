package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"posam/util/blockingqueue"
	"reflect"
	"runtime"
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

type InstructionMapt map[string]reflect.Type

var InstructionMap = make(InstructionMapt)

func (m InstructionMapt) Set(k string, v interface{}) {
	m[k] = reflect.TypeOf(v)
}

func (m InstructionMapt) Get(name string) (reflect.Value, error) {
	if t, ok := m[name]; ok {
		return reflect.New(t), nil
	}
	return reflect.Value{}, fmt.Errorf("Invalid instruction")
}

type Statement struct {
	StatementGroup  *StatementGroup
	InstructionName string
	Arguments       []string
	IgnoreError     bool
}

type StatementGroup struct {
	Execution ExecutionType
	Stack     *Stack
	ItemList  *blockingqueue.BlockingQueue
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

func InitParser(instructionMap InstructionMapt) error {
	for k, v := range instructionMap {
		InstructionMap[k] = v
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

func (s *Statement) Run(completec chan<- interface{}) Response {
	resp := Response{}
	if _, ok := InstructionMap[s.InstructionName]; !ok {
		resp.Error = fmt.Errorf("Invalid instruction %q", s.InstructionName)
	} else {
		instructionInstancev, err := InstructionMap.Get(s.InstructionName)
		reflect.Indirect(instructionInstancev).FieldByName("Env").Set(reflect.ValueOf(s.StatementGroup.Stack))

		arguments := make([]reflect.Value, len(s.Arguments))
		for i, _ := range s.Arguments {
			arguments[i] = reflect.ValueOf(s.Arguments[i])
		}

		outputValueList := instructionInstancev.MethodByName("Execute").Call(arguments)
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
		instructionInstancev.MethodByName("IssueError").Call([]reflect.Value{reflect.ValueOf(message)})
	}
	return resp
}

func (s *Statement) Execute(terminatec <-chan interface{}, completec chan<- interface{}) <-chan Response {

	respc := make(chan Response)
	go func() {
		defer close(respc)
		defer runtime.GC()
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
	}()
	log.Printf("'%s: %s' execute thread exits\n", s.InstructionName, s.Arguments)
	//time.Sleep(1 * time.Second)
	return respc
}

func tryRead(terminatec <-chan interface{}, inputc <-chan Response) <-chan Response {
	outputc := make(chan Response)
	go func() {
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
	}()
	return outputc
}

func bridge(terminatec <-chan interface{}, chanc <-chan <-chan Response) <-chan Response {
	valStream := make(chan Response)
	go func() {
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
			go func() {
				for val := range tryRead(terminatec, stream) {
					select {
					case valStream <- val:
					case <-terminatec:
					}
				}
			}()

		}
	}()
	return valStream
}

func (g *StatementGroup) ExecuteAsync(terminatec <-chan interface{}, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== ASYNC ====")
	respcc := make(chan (<-chan Response))

	go func() {
		var wg sync.WaitGroup
		wg.Add(g.ItemList.Length())
		defer close(respcc)

		for itemi := range g.ItemList.Iter() {
			itemInterface := itemi.Value
			completec := make(chan interface{})
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				go func() {
					respcc <- item.Execute(terminatec, completec)
				}()
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				go func() {
					respcc <- item.Execute(terminatec, completec)
				}()
			default:
				log.Printf("NO MATCH %T!\n", t)
			}

			go func() {
				defer wg.Done()
				<-completec
			}()
		}

		wg.Wait()
		pcompletec <- true
	}()

	return respcc
}

func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== SYNC ====")
	respcc := make(chan (<-chan Response))

	go func() {
		defer close(respcc)

		length := g.ItemList.Length()
		for i := 0; i < length; i++ {
			completec := make(chan interface{})
			itemInterface, err := g.ItemList.Get(i)
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
			varCur, _ := g.Stack.Get("SYS_CUR")
			cur64 := varCur.Value.(int64)
			cur := int(cur64)
			if cur > 0 {
				target := cur - 2
				if target == i {
					target += 1
				}
				i = target
				varCur.Value = int64(0)
			} else if cur < 0 {
				varCur.Value = int64(0)
				break
			}
		}

		pcompletec <- true
	}()

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
	stack *Stack,
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
		statementGroup.Stack = NewStack()
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		statement, _ := ParseLine(line)
		statement.StatementGroup = statementGroup
		switch statement.InstructionName {
		case "IMPORT":
			subStatementGroup, _ := ParseFile(
				NewStack(statementGroup.Stack),
				statement.Arguments[0],
				SYNC)
			statementGroup.ItemList.Append(subStatementGroup)
		case "ASYNC":
			subStatementGroup, _ := ParseFile(
				NewStack(statementGroup.Stack),
				statement.Arguments[0],
				ASYNC)
			statementGroup.ItemList.Append(subStatementGroup)
		default:
			statementGroup.ItemList.Append(statement)
		}
	}
	return statementGroup, nil

}
