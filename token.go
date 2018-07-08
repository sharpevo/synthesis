package commandparser

import (
	"fmt"
	"strconv"
	"time"
)

type Commander interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type Command struct {
	title string
	path  string
}

func (c Command) Title() string {
	return c.title
}

func (c *Command) SetTitle(title string) {
	c.title = title
}

func (c *Command) Execute(args ...string) (interface{}, error) {
	return "", nil
}

type CommandImport struct {
	Command
}

var Import CommandImport

type CommandAsync struct {
	Command
}

var Async CommandAsync

type CommandRetry struct {
	Command
}

var Retry CommandRetry

type CommandMove struct {
	Command
}

func (c *CommandMove) Execute(args ...string) (interface{}, error) {
	if c.isMovable() {
		result := fmt.Sprintf("Movable: %s", args[0])
		return result, nil
	} else {
		fmt.Println("Can not move")
		return c.Command.Execute(args...)
	}
}

func (c *CommandMove) isMovable() bool {
	return true
}

type CommandMoveX struct {
	CommandMove
}

var MoveX CommandMoveX

type CommandMoveY struct {
	CommandMove
}

var MoveY CommandMoveY

type CommandMoveZ struct {
	CommandMove
}

var MoveZ CommandMoveZ

type CommandSleep struct {
	Command
}

var Sleep CommandSleep

func (c *CommandSleep) Execute(args ...string) (interface{}, error) {
	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	return fmt.Sprintf("sleep %d seconds", seconds), nil
}

func Init() {
	Import.SetTitle("IMPORT")
	Async.SetTitle("ASYNC")
	Retry.SetTitle("RETRY")
	MoveX.SetTitle("MOVEX")
	MoveY.SetTitle("MOVEY")
	MoveZ.SetTitle("MOVEZ")
}
