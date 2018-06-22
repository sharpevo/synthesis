package commandparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ExecutionType int

const (
	SYNC ExecutionType = iota
	ASYNC
)

var CommandMap map[string]Commander

type Statement struct {
	CommandName string
	Arguments   []string
}

type StatementGroup struct {
	Execution ExecutionType
	ItemList  []interface{}
}

type Response struct {
	Error  error
	Output interface{}
}

type Info struct {
	Line   int
	Column int
}

func InitParser(commandMap map[string]Commander) error {
	CommandMap = commandMap
	return nil
}

func ParseLine(line string) (*Statement, error) {
	itemList := strings.Split(line, " ")
	statement := &Statement{}
	if len(itemList) < 2 {
		return statement, fmt.Errorf("Error: %s", "Invalid syntax")
	}
	statement.CommandName = itemList[0]
	statement.Arguments = itemList[1:]
	return statement, nil
}

func (s *Statement) Execute(terminatec <-chan interface{}) <-chan Response {
	fmt.Println("Run", s.CommandName)
	respc := make(chan Response)
	resp := Response{}
	go func() {
		defer close(respc)
		select {
		case <-terminatec:
			resp.Error = fmt.Errorf("Terminated %q", s.CommandName)
			respc <- resp
		default:
			if _, ok := CommandMap[s.CommandName]; !ok {
				resp.Error = fmt.Errorf("Invalid command %q", s.CommandName)
			} else {
				command := CommandMap[s.CommandName]
				output, _ := command.Execute(s.Arguments...)
				resp.Output = output
				fmt.Println("Run Output", s.CommandName, output)
			}
			respc <- resp
		}
		fmt.Println("Run Complete", s.CommandName, resp.Output)
	}()
	//time.Sleep(500 * time.Millisecond)
	return respc
}

func serialize(terminatec <-chan interface{}, respcs ...<-chan Response) <-chan Response {
	fmt.Println("Serialize", len(respcs))
	var wg sync.WaitGroup
	resultc := make(chan Response)

	swallow := func(respc <-chan Response) {
		fmt.Println("swallow")
		for resp := range respc {
			select {
			case <-terminatec:
				return
			case resultc <- resp:
				fmt.Printf("swallow assign %q\n", resp.Output)
			}
		}
		//time.Sleep(3 * time.Second)
		wg.Done()
		fmt.Println("swallow done")
	}

	wg.Add(len(respcs))
	for _, respc := range respcs {
		go swallow(respc)
	}

	go func() {
		wg.Wait()
		fmt.Println("DONE")
		//panic("done")
		//time.Sleep(3 * time.Second)
		close(resultc)
	}()
	fmt.Println("Finished")
	//time.Sleep(1 * time.Second)
	return resultc
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
			for val := range tryRead(terminatec, stream) {
				select {
				case valStream <- val:
				case <-terminatec:
				}
			}

		}
	}()
	return valStream
}

//func write(chanc <-chan <-chan Response) <-chan <-chan Response{
//go func(){
//}
//}

func (g *StatementGroup) ExecuteAsync(terminatec <-chan interface{}) <-chan <-chan Response {

	//respcList := []<-chan Response{}
	respcc := make(chan (<-chan Response))

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(g.ItemList))
		defer close(respcc)

		for _, itemInterface := range g.ItemList {
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				go func() {
					respcc <- item.Execute(terminatec)
					wg.Done()
					//respcList = append(respcList, item.Execute(terminatec))
				}()
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				go func() {
					respcc <- item.Execute(terminatec, nil)
					wg.Done()
					//respcList = append(respcList, item.Execute(terminatec, nil))
				}()
			default:
				fmt.Printf("NO MATCH %T!\n", t)
			}
		}
		wg.Wait()
	}()

	//return serialize(terminatec, respcList...)
	return respcc
}

//func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}) <-chan Response {
func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}) <-chan <-chan Response {
	//respcList := []<-chan Response{}

	respcc := make(chan (<-chan Response))

	go func() {
		defer close(respcc)

		for i := 0; i < len(g.ItemList); i++ {
			itemInterface := g.ItemList[i]
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				if item.CommandName == "RETRY" &&
					// TODO: previous error
					true {
					lineNum, _ := strconv.Atoi(item.Arguments[0])
					count, _ := strconv.Atoi(item.Arguments[1])
					if count < 1 {
						// TODO: panic
						fmt.Printf(
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
					respcc <- item.Execute(terminatec)
					//respcList = append(respcList, item.Execute(terminatec))
					fmt.Println("ExecuteSync Process", item)
				}
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				respcc <- item.Execute(terminatec, nil)
				//respcList = append(respcList, item.Execute(terminatec, nil))
			default:
				fmt.Printf("NO MATCH %T!\n", t)
			}
		}

	}()
	fmt.Println("ExecuteSync completed")
	time.Sleep(1 * time.Second)

	//return serialize(terminatec, respcList...)
	return respcc
}

//func (g *StatementGroup) Execute(terminatec <-chan interface{}, parentWg *sync.WaitGroup) <-chan Response {
func (g *StatementGroup) Execute(terminatec <-chan interface{}, parentWg *sync.WaitGroup) <-chan Response {

	//for {
	select {
	case <-terminatec:
		//respcc := make(chan (<-chan Response))
		//defer close(respcc)
		respc := make(chan Response)
		defer close(respc)
		resp := Response{}
		resp.Error = fmt.Errorf("Terminated %q", "Statement Group")
		respc <- resp
		return respc
		//respcc <- respc
		//return respcc
	default:
		switch g.Execution {
		case SYNC:
			return bridge(terminatec, g.ExecuteSync(terminatec))
		case ASYNC:
			return bridge(terminatec, g.ExecuteAsync(terminatec))
		}
	}
	//}
	resultc := make(chan Response)
	close(resultc)
	return resultc
	//resultcc := make(chan (<-chan Response))
	//close(resultcc)
	//return resultcc
}

//func (g *StatementGroup) xExecute(terminatec <-chan interface{}, parentWg *sync.WaitGroup) <-chan Response {
//respc := make(<-chan Response)

//go func() {
//defer close(respc)
//for {
//select {
//case <-terminatec:
//resp := Response{}
//resp.Error = fmt.Errorf("Terminated %q", "Statement Group")
//respc <- resp
//default:
//switch g.Execution {
//case SYNC:
//respc = g.ExecuteSync(terminatec)
////respc <- <-g.ExecuteSync(terminatec)
////return g.ExecuteSync(terminatec)
//case ASYNC:
//respc = g.ExecuteAsync(terminatec)
////return g.ExecuteAsync(terminatec)
////resp := <-g.ExecuteAsync(terminatec)
////respc <- resp
//}
//}
//}
//}()

//return respc
//}

func ParseFile(
	filePath string,
	execution ExecutionType) (*StatementGroup, error) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	statementGroup := StatementGroup{
		Execution: execution,
	}

	return ParseReader(file, &statementGroup)
}

func ParseReader(reader io.Reader, statementGroup *StatementGroup) (*StatementGroup, error) {

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		statement, _ := ParseLine(line)
		switch statement.CommandName {
		case "IMPORT":
			subStatementGroup, _ := ParseFile(
				statement.Arguments[0],
				SYNC)
			statementGroup.ItemList = append(
				statementGroup.ItemList,
				subStatementGroup)
		case "ASYNC":
			subStatementGroup, _ := ParseFile(
				statement.Arguments[0],
				ASYNC)
			statementGroup.ItemList = append(
				statementGroup.ItemList,
				subStatementGroup)
		default:
			statementGroup.ItemList = append(
				statementGroup.ItemList,
				statement)
		}
	}
	return statementGroup, nil

}
