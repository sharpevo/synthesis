package commandparser

import (
	"fmt"
	"strings"
)

type Statement struct {
	CommandName string
	Arguments   []string
}

func ParseLine(line string) (*Statement, error) {
	itemList := strings.Split(line, " ")
	statement := &Statement{}
	if len(itemList) > 1 {
		statement.CommandName = itemList[0]
		statement.Arguments = itemList[1:]
	} else {
		return statement, fmt.Errorf("Error: %s", "Invalid syntax")
	}
	return statement, nil
}

func (s *Statement) Execute() (string, error) {
	return CommandMap[s.CommandName](s.Arguments...)
}
