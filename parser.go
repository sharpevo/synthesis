package commandparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ExecutionType int

const (
	SYNC ExecutionType = iota
	ASYNC
)

type Statement struct {
	CommandName string
	Arguments   []string
}

type StatementGroup struct {
	Execution ExecutionType
	ItemList  []interface{}
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

func (s *Statement) Execute() (string, error) {
	return CommandMap[s.CommandName](s.Arguments...)
}

func (g *StatementGroup) Execute() ([]string, error) {
	resultList := []string{}
	switch g.Execution {
	case SYNC:
		for _, itemInterface := range g.ItemList {
			switch t := itemInterface.(type) {
			case Statement, *Statement:
				item, _ := itemInterface.(*Statement)
				result, _ := item.Execute()
				resultList = append(resultList, result)
			case StatementGroup, *StatementGroup:
				item, _ := itemInterface.(*StatementGroup)
				result, _ := item.Execute()
				resultList = append(resultList, result...)
			default:
				fmt.Printf("NO MATCH %T!\n", t)
			}
		}
	}
	return resultList, nil
}

func ParseFile(
	filePath string,
	execution ExecutionType) (*StatementGroup, error) {
	statementGroup := StatementGroup{
		Execution: execution,
	}

	file, err := os.Open(filePath)
	if err != nil {
		return &statementGroup, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
	return &statementGroup, nil
}
