package commandparser

import (
//"fmt"
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

func Init() {
	Import.SetTitle("IMPORT")
	Async.SetTitle("ASYNC")
	Retry.SetTitle("RETRY")
}
