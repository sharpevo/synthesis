package commandparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	//"time"
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
	respc := make(chan Response)
	resp := Response{}
	go func() {
		defer close(respc)
		for {
			select {
			case <-terminatec:
				resp.Error = fmt.Errorf("Terminated %q", s.CommandName)
				respc <- resp
			default:
				//fmt.Println(">>", s.CommandName)
				if _, ok := CommandMap[s.CommandName]; !ok {
					resp.Error = fmt.Errorf("Invalid command %q", s.CommandName)
				} else {
					command := CommandMap[s.CommandName]
					output, _ := command.Execute(s.Arguments...)
					resp.Output = output
				}
				respc <- resp
			}
		}
	}()
	//time.Sleep(1 * time.Second)
	return respc
}

func (g *StatementGroup) ExecuteAsync(terminatec <-chan interface{}) (outputList []string) {
	infoCh := make(chan Info)
	outputCh := make(chan interface{})
	errorCh := make(chan error)
	defer close(outputCh)
	defer close(errorCh)
	defer close(infoCh)
	var wg sync.WaitGroup
	wg.Add(len(g.ItemList))
	for _, itemInterface := range g.ItemList {
		switch t := itemInterface.(type) {
		case Statement, *Statement:
			item, _ := itemInterface.(*Statement)
			go func() {
				resp := <-item.Execute(terminatec)
				if resp.Error != nil {
					fmt.Println(resp.Error)
				}
				outputCh <- resp.Output
			}()
		case StatementGroup, *StatementGroup:
			item, _ := itemInterface.(*StatementGroup)
			wg.Add(len(item.ItemList) - 1)
			go func() {
				results, _ := item.Execute(terminatec, nil)
				for _, result := range results {
					outputCh <- result
				}
			}()
		default:
			fmt.Printf("NO MATCH %T!\n", t)
		}
	}
	go func() {
		for {
			select {
			case output, ok := <-outputCh:
				if ok {
					outputList = append(outputList, fmt.Sprintf("%v", output))
					wg.Done()
				}
			case err, ok := <-errorCh:
				if ok {
					// TODO: interruption
					fmt.Println(ok, err)
					return
				}
			case info, ok := <-infoCh:
				if ok {
					// TODO: location
					fmt.Println(info)
					return
				}
			}
		}
	}()
	wg.Wait()
	return
}

func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}) (outputList []string) {

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
				resp := <-item.Execute(terminatec)
				if resp.Error != nil {
					fmt.Println(resp.Error)
				}
				outputList = append(outputList, fmt.Sprintf("%v", resp.Output))
			}
		case StatementGroup, *StatementGroup:
			item, _ := itemInterface.(*StatementGroup)
			result, _ := item.Execute(terminatec, nil)
			outputList = append(outputList, result...)
		default:
			fmt.Printf("NO MATCH %T!\n", t)
		}
	}
	return
}

func (g *StatementGroup) Execute(terminatec <-chan interface{}, parentWg *sync.WaitGroup) (outputList []string, err error) {
	err = nil
	switch g.Execution {
	case SYNC:
		outputList = append(outputList, g.ExecuteSync(terminatec)...)
	case ASYNC:
		outputList = append(outputList, g.ExecuteAsync(terminatec)...)
	}
	if parentWg != nil {
		parentWg.Done()
	}

	return
}

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
