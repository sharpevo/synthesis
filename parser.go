package commandparser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

func (s *Statement) Run() Response {
	resp := Response{}
	if _, ok := CommandMap[s.CommandName]; !ok {
		resp.Error = fmt.Errorf("Invalid command %q", s.CommandName)
	} else {
		command := CommandMap[s.CommandName]
		output, err := command.Execute(s.Arguments...)
		resp.Output = output
		resp.Error = err
		log.Printf("'%s: %s' produces %q\n", s.CommandName, s.Arguments, output)
	}
	return resp
}

func (s *Statement) Execute(terminatec <-chan interface{}, suspended *bool, completec chan<- interface{}) <-chan Response {

	//fmt.Println("+", *suspended)
	if *suspended {
		log.Printf("'%s: %s' suspended\n", s.CommandName, s.Arguments)
		for {
			//fmt.Println("-", *suspended)
			if !*suspended {
				log.Printf("'%s: %s' resumed\n", s.CommandName, s.Arguments)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	respc := make(chan Response)
	go func() {

		defer close(respc)
		for {
			select {
			case <-terminatec:
				resp := Response{}
				log.Printf("Termiante '%s: %s'\n\n", s.CommandName, s.Arguments)
				resp.Error = fmt.Errorf("Terminated %q", s.CommandName)
				respc <- resp
				if completec != nil {
					completec <- true
				}
				return
			case respc <- s.Run():
				log.Printf("'%s: %s' done\n", s.CommandName, s.Arguments)
				if completec != nil {
					completec <- true
				}
				return
			}
		}
	}()
	log.Printf("'%s: %s' execute thread exits\n", s.CommandName, s.Arguments)
	//time.Sleep(1 * time.Second)
	return respc
}

func serialize(terminatec <-chan interface{}, respcs ...<-chan Response) <-chan Response {
	var wg sync.WaitGroup
	resultc := make(chan Response)

	swallow := func(respc <-chan Response) {
		for resp := range respc {
			select {
			case <-terminatec:
				return
			case resultc <- resp:
			}
		}
		wg.Done()
	}

	wg.Add(len(respcs))
	for _, respc := range respcs {
		go swallow(respc)
	}

	go func() {
		wg.Wait()
		close(resultc)
	}()
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

func (g *StatementGroup) ExecuteAsync(terminatec <-chan interface{}, suspended *bool, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== ASYNC ====")
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
					respcc <- item.Execute(terminatec, suspended, nil)
					wg.Done()
				}()
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				go func() {
					respcc <- item.Execute(terminatec, suspended, nil)
					wg.Done()
				}()
			default:
				log.Printf("NO MATCH %T!\n", t)
			}
		}
		wg.Wait()
		if pcompletec != nil {
			pcompletec <- true
		}
	}()

	return respcc
}

func (g *StatementGroup) ExecuteSync(terminatec <-chan interface{}, suspended *bool, pcompletec chan<- interface{}) <-chan <-chan Response {

	log.Println("==== SYNC ====")
	respcc := make(chan (<-chan Response))

	go func() {
		defer close(respcc)

		for i := 0; i < len(g.ItemList); i++ {
			completec := make(chan interface{})
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
					respcc <- item.Execute(terminatec, suspended, completec)
					//log.Printf("'%s: %s' complet", item.CommandName, item.Arguments)
				}

			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				respcc <- item.Execute(terminatec, suspended, completec)
			default:
				log.Printf("NO MATCH %T!\n", t)
			}

			<-completec
		}

		if pcompletec != nil {
			pcompletec <- true
		}
	}()

	return respcc
}

func (g *StatementGroup) Execute(terminatec <-chan interface{}, suspended *bool, completec chan<- interface{}) <-chan Response {
	switch g.Execution {
	case SYNC:
		return bridge(terminatec, g.ExecuteSync(terminatec, suspended, completec))
	case ASYNC:
		return bridge(terminatec, g.ExecuteAsync(terminatec, suspended, completec))
	}
	resultc := make(chan Response)
	close(resultc)
	return resultc
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
